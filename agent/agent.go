/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/orb-community/orb/agent/backend"
	"github.com/orb-community/orb/agent/cloud_config"
	"github.com/orb-community/orb/agent/config"
	manager "github.com/orb-community/orb/agent/policyMgr"
	"github.com/orb-community/orb/buildinfo"
	"go.uber.org/zap"
)

var (
	ErrMqttConnection = errors.New("failed to connect to a broker")
)

type Agent interface {
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
	Stop(ctx context.Context)
	RestartAll(ctx context.Context, reason string) error
	RestartBackend(ctx context.Context, backend string, reason string) error
}

type orbAgent struct {
	logger            *zap.Logger
	config            config.Config
	client            mqtt.Client
	agent_id          string
	db                *sqlx.DB
	backends          map[string]backend.Backend
	backendState      map[string]*backend.State
	cancelFunction    context.CancelFunc
	rpcFromCancelFunc context.CancelFunc

	asyncContext context.Context

	hbTicker        *time.Ticker
	heartbeatCtx    context.Context
	heartbeatCancel context.CancelFunc

	// Agent RPC channel, configured from command line
	baseTopic         string
	rpcToCoreTopic    string
	rpcFromCoreTopic  string
	capabilitiesTopic string
	heartbeatsTopic   string
	logTopic          string

	// Retry Mechanism to ensure the Request is received
	groupRequestTicker     *time.Ticker
	groupRequestSucceeded  context.CancelFunc
	policyRequestTicker    *time.Ticker
	policyRequestSucceeded context.CancelFunc

	// AgentGroup channels sent from core
	groupsInfos map[string]GroupInfo

	policyManager manager.PolicyManager
}

const retryRequestDuration = time.Second
const retryRequestFixedTime = 15
const retryDurationIncrPerAttempts = 10
const retryMaxAttempts = 4

type GroupInfo struct {
	Name      string
	ChannelID string
}

var _ Agent = (*orbAgent)(nil)

func New(logger *zap.Logger, c config.Config) (Agent, error) {
	logger.Info("using local config db", zap.String("filename", c.OrbAgent.DB.File))
	db, err := sqlx.Connect("sqlite3", c.OrbAgent.DB.File)
	if err != nil {
		return nil, err
	}

	pm, err := manager.New(logger, c, db)
	if err != nil {
		logger.Error("error during create policy manager, exiting", zap.Error(err))
		return nil, err
	}
	if pm.GetRepo() == nil {
		logger.Error("policy manager failed to get repository", zap.Error(err))
		return nil, err
	}
	return &orbAgent{logger: logger, config: c, policyManager: pm, db: db, groupsInfos: make(map[string]GroupInfo)}, nil
}

func (a *orbAgent) startBackends(agentCtx context.Context) error {
	a.logger.Info("registered backends", zap.Strings("values", backend.GetList()))
	a.logger.Info("requested backends", zap.Any("values", a.config.OrbAgent.Backends))
	if len(a.config.OrbAgent.Backends) == 0 {
		return errors.New("no backends specified")
	}
	a.backends = make(map[string]backend.Backend, len(a.config.OrbAgent.Backends))
	a.backendState = make(map[string]*backend.State)
	for name, configurationEntry := range a.config.OrbAgent.Backends {
		if !backend.HaveBackend(name) {
			return errors.New("specified backend does not exist: " + name)
		}
		be := backend.GetBackend(name)
		configuration := structs.Map(a.config.OrbAgent.Otel)
		configuration["agent_tags"] = a.config.OrbAgent.Tags
		if err := be.Configure(a.logger, a.policyManager.GetRepo(), configurationEntry, configuration); err != nil {
			a.logger.Info("failed to configure backend", zap.String("backend", name), zap.Error(err))
			return err
		}
		backendCtx := context.WithValue(agentCtx, "routine", name)
		if a.config.OrbAgent.Cloud.MQTT.Id != "" {
			backendCtx = context.WithValue(backendCtx, "agent_id", a.config.OrbAgent.Cloud.MQTT.Id)
		} else {
			backendCtx = context.WithValue(backendCtx, "agent_id", "auto-provisioning-without-id")
		}
		a.backends[name] = be
		initialState := be.GetInitialState()
		a.backendState[name] = &backend.State{
			Status:        initialState,
			LastRestartTS: time.Now(),
		}
		if err := be.Start(context.WithCancel(backendCtx)); err != nil {
			a.logger.Info("failed to start backend", zap.String("backend", name), zap.Error(err))
			var errMessage string
			if initialState == backend.BackendError {
				errMessage = err.Error()
			}
			a.backendState[name] = &backend.State{
				Status:        initialState,
				LastError:     errMessage,
				LastRestartTS: time.Now(),
			}
			return err
		}
	}
	return nil
}

func (a *orbAgent) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	startTime := time.Now()
	defer func(t time.Time) {
		a.logger.Debug("Startup of agent execution duration", zap.String("Start() execution duration", time.Since(t).String()))
	}(startTime)
	agentCtx := context.WithValue(ctx, "routine", "agentRoutine")
	asyncCtx, cancelAllAsync := context.WithCancel(context.WithValue(ctx, "routine", "asyncParent"))
	a.asyncContext = asyncCtx
	a.rpcFromCancelFunc = cancelAllAsync
	a.cancelFunction = cancelFunc
	a.logger.Info("agent started", zap.String("version", buildinfo.GetVersion()), zap.Any("routine", agentCtx.Value("routine")))
	mqtt.CRITICAL = &agentLoggerCritical{a: a}
	mqtt.ERROR = &agentLoggerError{a: a}

	if a.config.OrbAgent.Debug.Enable {
		a.logger.Info("debug logging enabled")
		mqtt.DEBUG = &agentLoggerDebug{a: a}
	}

	ccm, err := cloud_config.New(a.logger, a.config, a.db)
	if err != nil {
		return err
	}
	cloudConfig, err := ccm.GetCloudConfig()
	if err != nil {
		return err
	}

	commsCtx := context.WithValue(agentCtx, "routine", "comms")
	if err := a.startComms(commsCtx, cloudConfig); err != nil {
		a.logger.Error("could not start mqtt client")
		return err
	}

	if err := a.startBackends(ctx); err != nil {
		return err
	}

	a.logonWithHeartbeat()

	return nil
}

func (a *orbAgent) logonWithHeartbeat() {
	a.hbTicker = time.NewTicker(HeartbeatFreq)
	a.heartbeatCtx, a.heartbeatCancel = a.extendContext("heartbeat")
	go a.sendHeartbeats(a.heartbeatCtx, a.heartbeatCancel)
	a.logger.Info("heartbeat routine started")
}

func (a *orbAgent) logoffWithHeartbeat(ctx context.Context) {
	a.logger.Debug("stopping heartbeat, going offline status", zap.Any("routine", ctx.Value("routine")))
	if a.heartbeatCtx != nil {
		a.heartbeatCancel()
	}
	if a.client != nil && a.client.IsConnected() {
		a.unsubscribeGroupChannels()
		if token := a.client.Unsubscribe(a.rpcFromCoreTopic); token.Wait() && token.Error() != nil {
			a.logger.Warn("failed to unsubscribe to RPC channel", zap.Error(token.Error()))
		}
	}
}
func (a *orbAgent) Stop(ctx context.Context) {
	a.logger.Info("routine call for stop agent", zap.Any("routine", ctx.Value("routine")))
	if a.rpcFromCancelFunc != nil {
		a.rpcFromCancelFunc()
	}
	for name, b := range a.backends {
		if state, _, _ := b.GetRunningStatus(); state == backend.Running {
			a.logger.Debug("stopping backend", zap.String("backend", name))
			if err := b.Stop(ctx); err != nil {
				a.logger.Error("error while stopping the backend", zap.String("backend", name))
			}
		}
	}
	a.logoffWithHeartbeat(ctx)
	if a.client != nil && a.client.IsConnected() {
		a.client.Disconnect(0)
	}
	a.logger.Debug("stopping agent with number of go routines and go calls", zap.Int("goroutines", runtime.NumGoroutine()), zap.Int64("gocalls", runtime.NumCgoCall()))
	if a.policyRequestSucceeded != nil {
		a.policyRequestSucceeded()
	}
	if a.groupRequestSucceeded != nil {
		a.groupRequestSucceeded()
	}
	defer a.cancelFunction()
}

func (a *orbAgent) RestartBackend(ctx context.Context, name string, reason string) error {
	if !backend.HaveBackend(name) {
		return errors.New("specified backend does not exist: " + name)
	}

	be := a.backends[name]
	a.logger.Info("restarting backend", zap.String("backend", name), zap.String("reason", reason))
	a.backendState[name].RestartCount += 1
	a.backendState[name].LastRestartTS = time.Now()
	a.backendState[name].LastRestartReason = reason
	a.logger.Info("removing policies", zap.String("backend", name))
	if err := a.policyManager.RemoveBackendPolicies(be, true); err != nil {
		a.logger.Error("failed to remove policies", zap.String("backend", name), zap.Error(err))
	}
	configuration := structs.Map(a.config.OrbAgent.Otel)
	configuration["agent_tags"] = a.config.OrbAgent.Tags
	if err := be.Configure(a.logger, a.policyManager.GetRepo(), a.config.OrbAgent.Backends[name], configuration); err != nil {
		return err
	}
	a.logger.Info("resetting backend", zap.String("backend", name))

	if err := be.FullReset(ctx); err != nil {
		a.backendState[name].LastError = fmt.Sprintf("failed to reset backend: %v", err)
		a.logger.Error("failed to reset backend", zap.String("backend", name), zap.Error(err))
	}
	be.SetCommsClient(a.agent_id, &a.client, fmt.Sprintf("%s/?/%s", a.baseTopic, name))

	if err := a.sendAgentPoliciesReq(); err != nil {
		a.logger.Error("failed to send agent policies request", zap.Error(err))
	}
	return nil
}

func (a *orbAgent) restartComms(ctx context.Context) error {
	if a.client != nil && a.client.IsConnected() {
		a.unsubscribeGroupChannels()
	}
	ccm, err := cloud_config.New(a.logger, a.config, a.db)
	if err != nil {
		return err
	}
	cloudConfig, err := ccm.GetCloudConfig()
	if err != nil {
		return err
	}
	if err := a.startComms(ctx, cloudConfig); err != nil {
		a.logger.Error("could not restart mqtt client")
		return err
	}
	return nil
}

func (a *orbAgent) RestartAll(ctx context.Context, reason string) error {
	if a.config.OrbAgent.Cloud.MQTT.Id != "" {
		ctx = context.WithValue(ctx, "agent_id", a.config.OrbAgent.Cloud.MQTT.Id)
	} else {
		ctx = context.WithValue(ctx, "agent_id", "auto-provisioning-without-id")
	}
	a.logoffWithHeartbeat(ctx)
	a.logger.Info("restarting comms", zap.String("reason", reason))
	if err := a.restartComms(ctx); err != nil {
		a.logger.Error("failed to restart comms", zap.Error(err))
	}
	for name := range a.backends {
		a.logger.Info("restarting backend", zap.String("backend", name), zap.String("reason", reason))
		err := a.RestartBackend(ctx, name, reason)
		if err != nil {
			a.logger.Error("failed to restart backend", zap.Error(err))
		}
	}
	a.logger.Info("all backends and comms were restarted")

	return nil
}

func (a *orbAgent) extendContext(routine string) (context.Context, context.CancelFunc) {
	uuidTraceId := uuid.NewString()
	a.logger.Debug("creating context for receiving message", zap.String("routine", routine), zap.String("trace-id", uuidTraceId))
	return context.WithCancel(context.WithValue(context.WithValue(a.asyncContext, "routine", routine), "trace-id", uuidTraceId))
}
