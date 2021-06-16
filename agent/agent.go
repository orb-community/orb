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
	"github.com/ns1labs/orb"
	"github.com/ns1labs/orb/agent/backend"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

type Agent interface {
	Start() error
	Stop()
}

const HeartbeatFreq = 60 * time.Second

type orbAgent struct {
	logger   *zap.Logger
	config   Config
	client   mqtt.Client
	backends map[string]backend.Backend

	hbTicker *time.Ticker
	hbDone   chan bool

	rpcToCoreChannel    string
	rpcFromCoreChannel  string
	capabilitiesChannel string
	heartbeatsChannel   string
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
	a.rpcToCoreChannel = fmt.Sprintf("%s/%s", base, fleet.RPCToCoreChannel)
	a.rpcFromCoreChannel = fmt.Sprintf("%s/%s", base, fleet.RPCFromCoreChannel)
	a.capabilitiesChannel = fmt.Sprintf("%s/%s", base, fleet.CapabilitiesChannel)
	a.heartbeatsChannel = fmt.Sprintf("%s/%s", base, fleet.HeartbeatsChannel)
	a.logChannel = fmt.Sprintf("%s/%s", base, fleet.LogChannel)

}
func (a *orbAgent) sendSingleHeartbeat(t time.Time) {

	a.logger.Debug("heartbeat")

	hbData := make(map[string]interface{})
	hbData["ts"] = t.Unix()

	body, err := json.Marshal(hbData)
	if err != nil {
		a.logger.Error("error creating heartbeat data", zap.Error(err))
		return
	}

	if token := a.client.Publish(a.heartbeatsChannel, 1, false, body); token.Wait() && token.Error() != nil {
		a.logger.Error("error sending heartbeat", zap.Error(err))
	}
}

func (a *orbAgent) sendHeartbeats() {
	a.sendSingleHeartbeat(time.Now())
	for {
		select {
		case <-a.hbDone:
			return
		case t := <-a.hbTicker.C:
			a.sendSingleHeartbeat(t)
		}
	}
}

func (a *orbAgent) startComms() error {
	var err error
	a.client, err = a.connect()
	if err != nil {
		a.logger.Error("connection failed", zap.Error(err))
		return err
	}

	a.nameChannels()

	if token := a.client.Subscribe(a.rpcFromCoreChannel, 1, a.handleRPCFromCore); token.Wait() && token.Error() != nil {
		a.logger.Error("failed to subscribe to RPC channel", zap.Error(err))
		return err
	}

	err = a.sendCapabilities()
	if err != nil {
		a.logger.Error("failed to send agent capabilities", zap.Error(err))
		return err
	}

	a.hbTicker = time.NewTicker(HeartbeatFreq)
	a.hbDone = make(chan bool)
	go a.sendHeartbeats()

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
		Version: orb.GetVersion(),
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
	a.hbTicker.Stop()
	a.hbDone <- true
	if token := a.client.Unsubscribe(a.rpcFromCoreChannel); token.Wait() && token.Error() != nil {
		a.logger.Warn("failed to unsubscribe to RPC channel", zap.Error(token.Error()))
	}
	for _, be := range a.backends {
		if err := be.Stop(); err != nil {
			a.logger.Error("backend error while stopping", zap.Error(err))
		}
	}
	a.client.Disconnect(250)
}
