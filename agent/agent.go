/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"errors"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/cloud_config"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/agent/policyMgr"
	"github.com/ns1labs/orb/buildinfo"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

var (
	ErrMqttConnection = errors.New("failed to connect to a broker")
)

type Agent interface {
	Start() error
	Stop()
	RestartAll(reason string) error
	RestartBackend(backend string, reason string) error
}

type orbAgent struct {
	logger   *zap.Logger
	config   config.Config
	client   mqtt.Client
	db       *sqlx.DB
	backends map[string]backend.Backend

	hbTicker *time.Ticker
	hbDone   chan bool

	// Agent RPC channel, configured from command line
	baseTopic         string
	rpcToCoreTopic    string
	rpcFromCoreTopic  string
	capabilitiesTopic string
	heartbeatsTopic   string
	logTopic          string

	// Retry Mechanism to ensure the Request is received
	groupRequestTicker     *time.Ticker
	groupRequestSucceeded  chan bool
	policyRequestTicker    *time.Ticker
	policyRequestSucceeded chan bool

	// AgentGroup channels sent from core
	groupsInfos map[string]GroupInfo

	policyManager manager.PolicyManager
}

const retryRequestDuration = time.Second
const retryRequestFixedTime = 5
const retryDurationIncrPerAttempts = 10
const retryMaxAttempts = 5

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
		return nil, err
	}
	return &orbAgent{logger: logger, config: c, policyManager: pm, db: db, groupsInfos: make(map[string]GroupInfo)}, nil
}

func (a *orbAgent) startBackends() error {
	a.logger.Info("registered backends", zap.Strings("values", backend.GetList()))
	a.logger.Info("requested backends", zap.Any("values", a.config.OrbAgent.Backends))
	if len(a.config.OrbAgent.Backends) == 0 {
		return errors.New("no backends specified")
	}
	a.backends = make(map[string]backend.Backend, len(a.config.OrbAgent.Backends))
	for name, config := range a.config.OrbAgent.Backends {
		if !backend.HaveBackend(name) {
			return errors.New("specified backend does not exist: " + name)
		}
		be := backend.GetBackend(name)
		if err := be.Configure(a.logger, a.policyManager.GetRepo(), config, structs.Map(a.config.OrbAgent.Otel)); err != nil {
			return err
		}
		if err := be.Start(); err != nil {
			return err
		}
		a.backends[name] = be
	}
	return nil
}

func (a *orbAgent) Start() error {

	a.logger.Info("agent started", zap.String("version", buildinfo.GetVersion()))

	mqtt.CRITICAL = &agentLoggerCritical{a: a}
	mqtt.ERROR = &agentLoggerError{a: a}

	if a.config.OrbAgent.Debug.Enable {
		a.logger.Info("debug logging enabled")
		mqtt.DEBUG = &agentLoggerDebug{a: a}
	}

	if err := a.startBackends(); err != nil {
		return err
	}

	ccm, err := cloud_config.New(a.logger, a.config, a.db)
	if err != nil {
		return err
	}
	cloudConfig, err := ccm.GetCloudConfig()
	if err != nil {
		return err
	}

	if err := a.startComms(cloudConfig); err != nil {
		a.logger.Error("could not restart mqtt client")
		return err
	}

	a.hbTicker = time.NewTicker(HeartbeatFreq)
	a.hbDone = make(chan bool)
	a.groupRequestSucceeded = make(chan bool)
	a.policyRequestSucceeded = make(chan bool)
	go a.sendHeartbeats()

	return nil
}

func (a *orbAgent) Stop() {
	a.logger.Info("stopping agent")
	a.hbTicker.Stop()
	a.hbDone <- true
	a.closeRequestTickers()
	a.sendSingleHeartbeat(time.Now(), fleet.Offline)
	if token := a.client.Unsubscribe(a.rpcFromCoreTopic); token.Wait() && token.Error() != nil {
		a.logger.Warn("failed to unsubscribe to RPC channel", zap.Error(token.Error()))
	}
	a.unsubscribeGroupChannels()
	for _, be := range a.backends {
		if err := be.Stop(); err != nil {
			a.logger.Error("backend error while stopping", zap.Error(err))
		}
	}
	a.client.Disconnect(250)
}

func (a *orbAgent) closeRequestTickers() {
	a.groupRequestTicker.Stop()
	a.groupRequestSucceeded <- true
	a.policyRequestTicker.Stop()
	a.policyRequestSucceeded <- true
}

func (a *orbAgent) RestartBackend(name string, reason string) error {
	if !backend.HaveBackend(name) {
		return errors.New("specified backend does not exist: " + name)
	}

	be := a.backends[name]
	a.logger.Info("restarting backend", zap.String("backend", name), zap.String("reason", reason))
	a.logger.Info("removing policies", zap.String("backend", name))
	if err := a.policyManager.RemoveBackendPolicies(be, false); err != nil {
		a.logger.Error("failed to remove policies", zap.String("backend", name), zap.Error(err))
	}
	a.logger.Info("resetting backend", zap.String("backend", name))
	if err := be.FullReset(); err != nil {
		a.logger.Error("failed to reset backend", zap.String("backend", name), zap.Error(err))
	}
	a.logger.Info("reapplying policies", zap.String("backend", name))
	if err := a.policyManager.ApplyBackendPolicies(be); err != nil {
		a.logger.Error("failed to reapply policies", zap.String("backend", name), zap.Error(err))
	}
	return nil
}

func (a *orbAgent) restartComms() error {
	ccm, err := cloud_config.New(a.logger, a.config, a.db)
	if err != nil {
		return err
	}
	cloudConfig, err := ccm.GetCloudConfig()
	if err != nil {
		return err
	}
	a.closeRequestTickers()
	a.logger.Debug("restarting mqtt client")
	if err := a.startComms(cloudConfig); err != nil {
		a.logger.Error("could not restart mqtt client")
		return err
	}
	return nil
}

func (a *orbAgent) RestartAll(reason string) error {
	a.logger.Info("restarting comms")
	err := a.restartComms()
	if err != nil {
		a.logger.Error("failed to restart comms", zap.Error(err))
	}

	a.logger.Info("restarting all backends", zap.String("reason", reason))
	for name := range a.backends {
		a.logger.Info("restarting backend", zap.String("backend", name), zap.String("reason", reason))
		err = a.RestartBackend(name, reason)
		if err != nil {
			a.logger.Error("failed to restart backend", zap.Error(err))
		}
	}
	a.logger.Info("all backends and comms were restarted")

	return nil
}
