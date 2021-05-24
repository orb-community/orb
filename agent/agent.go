/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"time"
)

type Agent struct {
	logger *zap.Logger
	config Config

	rpcChannelToCore   string
	rpcChannelFromCore string
	client             mqtt.Client
}

func New(logger *zap.Logger, c Config) (*Agent, error) {
	return &Agent{logger: logger, config: c}, nil
}

func (a *Agent) connect() (mqtt.Client, error) {

	opts := mqtt.NewClientOptions().AddBroker(a.config.OrbAgent.MQTT["address"]).SetClientID(a.config.OrbAgent.MQTT["id"])
	opts.SetUsername(a.config.OrbAgent.MQTT["id"])
	opts.SetPassword(a.config.OrbAgent.MQTT["key"])
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		a.logger.Info("message on unknown channel", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))
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

func (a *Agent) Start() error {

	a.logger.Info("agent started")

	if a.config.Debug {
		mqtt.DEBUG = &agentLoggerDebug{a: a}
		a.logger.Debug("config", zap.Any("values", a.config))
	}
	mqtt.WARN = &agentLoggerWarn{a: a}
	mqtt.CRITICAL = &agentLoggerCritical{a: a}
	mqtt.ERROR = &agentLoggerError{a: a}

	var err error
	a.client, err = a.connect()
	if err != nil {
		a.logger.Error("connection failed", zap.Error(err))
		return err
	}

	a.rpcChannelToCore = fmt.Sprintf("channels/%s/messages/out", a.config.OrbAgent.MQTT["channel_id"])
	a.rpcChannelFromCore = fmt.Sprintf("channels/%s/messages/in", a.config.OrbAgent.MQTT["channel_id"])

	if token := a.client.Subscribe(a.rpcChannelFromCore, 1, a.handleRPC); token.Wait() && token.Error() != nil {
		a.logger.Error("failed to subscribe to RPC channel", zap.Error(err))
		return err
	}

	err = a.sendAgentInfo()
	if err != nil {
		a.logger.Error("failed to send AGENTINFO", zap.Error(err))
		return err
	}

	return nil
}

func (a *Agent) sendAgentInfo() error {

	agentInfo := make(map[string]string)
	agentInfo["version"] = "1.0"

	body, err := json.Marshal(agentInfo)
	if err != nil {
		return err
	}

	if token := a.client.Publish(a.rpcChannelToCore, 0, false, body); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (a *Agent) handleRPC(client mqtt.Client, message mqtt.Message) {
	a.logger.Info("RPC message", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))
}

func (a *Agent) Stop() {
	a.logger.Info("stopping agent")
	if token := a.client.Unsubscribe(a.rpcChannelFromCore); token.Wait() && token.Error() != nil {
		a.logger.Warn("failed to unsubscribe to RPC channel", zap.Error(token.Error()))
	}
	a.client.Disconnect(250)
}
