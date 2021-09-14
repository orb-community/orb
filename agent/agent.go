/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"errors"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/agent/cloud_config"
	"github.com/ns1labs/orb/agent/config"
	"github.com/ns1labs/orb/agent/policies"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

type Agent interface {
	Start() error
	Stop()
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
	rpcToCoreTopic    string
	rpcFromCoreTopic  string
	capabilitiesTopic string
	heartbeatsTopic   string
	logTopic          string

	// AgentGroup channels sent from core
	groupChannels []string

	policyManager policies.PolicyManager
}

var _ Agent = (*orbAgent)(nil)

func New(logger *zap.Logger, c config.Config) (Agent, error) {
	logger.Info("using local config db", zap.String("filename", c.OrbAgent.DB.File))
	db, err := sqlx.Connect("sqlite3", c.OrbAgent.DB.File)
	if err != nil {
		return nil, err
	}

	pm, err := policies.New(logger, c, db)
	if err != nil {
		return nil, err
	}
	return &orbAgent{logger: logger, config: c, policyManager: pm, db: db}, nil
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
		if err := be.Configure(a.logger, config); err != nil {
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

	a.logger.Info("agent started")

	mqtt.CRITICAL = &agentLoggerCritical{a: a}
	mqtt.ERROR = &agentLoggerError{a: a}

	ccm, err := cloud_config.New(a.logger, a.config, a.db)
	if err != nil {
		return err
	}
	cloudConfig, err := ccm.GetCloudConfig()
	if err != nil {
		return err
	}

	if err := a.startBackends(); err != nil {
		return err
	}

	if err := a.startComms(cloudConfig); err != nil {
		return err
	}

	return nil
}

func (a *orbAgent) Stop() {
	a.logger.Info("stopping agent")
	a.hbTicker.Stop()
	a.hbDone <- true
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
