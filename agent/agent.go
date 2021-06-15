/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/pkg/version"
	"go.uber.org/zap"
	"time"
)

type Agent interface {
	Start() error
	Stop()
}

type orbAgent struct {
	logger   *zap.Logger
	config   Config
	client   mqtt.Client
	backends map[string]backend.Backend

	rpcChannelToCore    string
	rpcChannelFromCore  string
	capabilitiesChannel string
	heartChannel        string
	logChannel          string
}

var _ Agent = (*orbAgent)(nil)

func New(logger *zap.Logger, c Config) (Agent, error) {
	return &orbAgent{logger: logger, config: c}, nil
}

func (a *orbAgent) connect() (mqtt.Client, error) {

	opts := mqtt.NewClientOptions().AddBroker(a.config.OrbAgent.MQTT["address"]).SetClientID(a.config.OrbAgent.MQTT["id"])
	opts.SetUsername(a.config.OrbAgent.MQTT["id"])
	opts.SetPassword(a.config.OrbAgent.MQTT["key"])
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		a.logger.Info("message on unknown channel, ignoring", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))
	})
	opts.SetPingTimeout(1 * time.Second)

	if !a.config.OrbAgent.TLS.Verify {
		opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return c, nil
}

func (a *orbAgent) nameChannels() {

	base := fmt.Sprintf("channels/%s/messages", a.config.OrbAgent.MQTT["channel_id"])
	a.rpcChannelToCore = fmt.Sprintf("%s/out", base)
	a.rpcChannelFromCore = fmt.Sprintf("%s/in", base)
	a.capabilitiesChannel = fmt.Sprintf("%s/agent", base)
	a.heartChannel = fmt.Sprintf("%s/hb", base)
	a.logChannel = fmt.Sprintf("%s/log", base)

}

func (a *orbAgent) startComms() error {
	var err error
	a.client, err = a.connect()
	if err != nil {
		a.logger.Error("connection failed", zap.Error(err))
		return err
	}

	a.nameChannels()

	if token := a.client.Subscribe(a.rpcChannelFromCore, 1, a.handleRPCFromCore); token.Wait() && token.Error() != nil {
		a.logger.Error("failed to subscribe to RPC channel", zap.Error(err))
		return err
	}

	err = a.sendCapabilities()
	if err != nil {
		a.logger.Error("failed to send agent info", zap.Error(err))
		return err
	}

	return nil
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

	if a.config.Debug {
		//mqtt.DEBUG = &agentLoggerDebug{a: a}
		a.logger.Debug("config", zap.Any("values", a.config))
	}
	mqtt.WARN = &agentLoggerWarn{a: a}
	mqtt.CRITICAL = &agentLoggerCritical{a: a}
	mqtt.ERROR = &agentLoggerError{a: a}

	if err := a.startBackends(); err != nil {
		return err
	}

	if err := a.startComms(); err != nil {
		return err
	}

	return nil
}

func (a *orbAgent) sendCapabilities() error {

	capabilities := make(map[string]interface{})

	capabilities["orb_agent"] = &OrbAgentInfo{
		Version: version.Version,
	}
	capabilities["backends"] = make(map[string]BackendInfo)
	for name, be := range a.backends {
		ver, err := be.Version()
		if err != nil {
			a.logger.Error("backend failed to retrieve version", zap.String("backend", name), zap.Error(err))
			continue
		}
		capabilities["backends"].(map[string]BackendInfo)[name] = BackendInfo{
			Version: ver,
		}
	}

	body, err := json.Marshal(capabilities)
	if err != nil {
		return err
	}

	if token := a.client.Publish(a.capabilitiesChannel, 1, false, body); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (a *orbAgent) handleRPCFromCore(client mqtt.Client, message mqtt.Message) {
	a.logger.Info("RPC message", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))
}

func (a *orbAgent) Stop() {
	a.logger.Info("stopping agent")
	if token := a.client.Unsubscribe(a.rpcChannelFromCore); token.Wait() && token.Error() != nil {
		a.logger.Warn("failed to unsubscribe to RPC channel", zap.Error(token.Error()))
	}
	for _, be := range a.backends {
		if err := be.Stop(); err != nil {
			a.logger.Error("backend error while stopping", zap.Error(err))
		}
	}
	a.client.Disconnect(250)
}
