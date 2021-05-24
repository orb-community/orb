/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"crypto/tls"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"time"
)

type Agent struct {
	logger *zap.Logger
	config Config

	rpcChannel string
	client     mqtt.Client
}

func New(logger *zap.Logger, c Config) (*Agent, error) {
	return &Agent{logger: logger, config: c}, nil
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (a *Agent) connect() (mqtt.Client, error) {

	opts := mqtt.NewClientOptions().AddBroker(a.config.OrbAgent.MQTT["address"]).SetClientID(a.config.OrbAgent.MQTT["id"])
	opts.SetUsername(a.config.OrbAgent.MQTT["id"])
	opts.SetPassword(a.config.OrbAgent.MQTT["key"])
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	// todo
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return c, nil
}

func (a *Agent) Start() error {

	a.logger.Info("agent started")

	mqtt.DEBUG = &agentLoggerDebug{a: a}
	mqtt.WARN = &agentLoggerWarn{a: a}
	mqtt.CRITICAL = &agentLoggerCritical{a: a}
	mqtt.ERROR = &agentLoggerError{a: a}

	var err error
	a.client, err = a.connect()
	if err != nil {
		a.logger.Error("connection failed", zap.Error(err))
		return err
	}

	a.rpcChannel = fmt.Sprintf("channels/%s/messages", a.config.OrbAgent.MQTT["channel_id"])

	if token := a.client.Subscribe(a.rpcChannel, 0, nil); token.Wait() && token.Error() != nil {
		a.logger.Error("failed to subscribe to RPC channel", zap.Error(err))
		return err
	}

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := a.client.Publish(a.rpcChannel, 0, false, text)
		token.Wait()
	}

	return nil
}

func (a *Agent) Stop() {
	a.logger.Info("stopping agent")
	if token := a.client.Unsubscribe(a.rpcChannel); token.Wait() && token.Error() != nil {
		a.logger.Warn("failed to unsubscribe to RPC channel", zap.Error(token.Error()))
	}
	a.client.Disconnect(250)
}
