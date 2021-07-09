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

	// Agent RPC channel, configured from command line
	rpcToCoreTopic    string
	rpcFromCoreTopic  string
	capabilitiesTopic string
	heartbeatsTopic   string
	logTopic          string

	// AgentGroup channels sent from core
	groupChannels []string
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

func (a *orbAgent) nameAgentRPCTopics() {

	base := fmt.Sprintf("channels/%s/messages", a.config.OrbAgent.MQTT["channel_id"])
	a.rpcToCoreTopic = fmt.Sprintf("%s/%s", base, fleet.RPCToCoreTopic)
	a.rpcFromCoreTopic = fmt.Sprintf("%s/%s", base, fleet.RPCFromCoreTopic)
	a.capabilitiesTopic = fmt.Sprintf("%s/%s", base, fleet.CapabilitiesTopic)
	a.heartbeatsTopic = fmt.Sprintf("%s/%s", base, fleet.HeartbeatsTopic)
	a.logTopic = fmt.Sprintf("%s/%s", base, fleet.LogTopic)

}
func (a *orbAgent) sendSingleHeartbeat(t time.Time, state fleet.State) {

	a.logger.Debug("heartbeat")

	hbData := fleet.Heartbeat{
		SchemaVersion: fleet.CurrentHeartbeatSchemaVersion,
		TimeStamp:     t,
		State:         state,
	}

	body, err := json.Marshal(hbData)
	if err != nil {
		a.logger.Error("error creating heartbeat data", zap.Error(err))
		return
	}

	if token := a.client.Publish(a.heartbeatsTopic, 1, false, body); token.Wait() && token.Error() != nil {
		a.logger.Error("error sending heartbeat", zap.Error(token.Error()))
	}
}

func (a *orbAgent) sendHeartbeats() {
	a.sendSingleHeartbeat(time.Now(), fleet.Online)
	for {
		select {
		case <-a.hbDone:
			return
		case t := <-a.hbTicker.C:
			a.sendSingleHeartbeat(t, fleet.Online)
		}
	}
}

func (a *orbAgent) unsubscribeGroupChannels() {
	for _, channel := range a.groupChannels {
		if token := a.client.Unsubscribe(channel); token.Wait() && token.Error() != nil {
			a.logger.Warn("failed to unsubscribe to group channel", zap.String("topic", channel), zap.Error(token.Error()))
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

	a.nameAgentRPCTopics()

	if token := a.client.Subscribe(a.rpcFromCoreTopic, 1, a.handleRPCFromCore); token.Wait() && token.Error() != nil {
		a.logger.Error("failed to subscribe to RPC topic", zap.String("topic", a.rpcFromCoreTopic), zap.Error(token.Error()))
		return token.Error()
	}

	err = a.sendCapabilities()
	if err != nil {
		a.logger.Error("failed to send agent capabilities", zap.Error(err))
		return err
	}

	err = a.sendGroupMembershipReq()
	if err != nil {
		a.logger.Error("failed to send group membership request", zap.Error(err))
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
	//	mqtt.WARN = &agentLoggerWarn{a: a}
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

	capabilities := fleet.Capabilities{
		SchemaVersion: fleet.CurrentCapabilitiesSchemaVersion,
		AgentTags:     a.config.OrbAgent.Tags,
		OrbAgent: fleet.OrbAgentInfo{
			Version: orb.GetVersion(),
		},
	}

	capabilities.Backends = make(map[string]fleet.BackendInfo)
	for name, be := range a.backends {
		ver, err := be.Version()
		if err != nil {
			a.logger.Error("backend failed to retrieve version", zap.String("backend", name), zap.Error(err))
			continue
		}
		capabilities.Backends[name] = fleet.BackendInfo{
			Version: ver,
		}
	}

	body, err := json.Marshal(capabilities)
	if err != nil {
		return err
	}

	if token := a.client.Publish(a.capabilitiesTopic, 1, false, body); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (a *orbAgent) sendGroupMembershipReq() error {

	payload := fleet.GroupMembershipReqRPCPayload{}

	data := fleet.RPC{
		SchemaVersion: fleet.CurrentRPCSchemaVersion,
		Func:          fleet.GroupMembershipReqRPCFunc,
		Payload:       payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if token := a.client.Publish(a.rpcToCoreTopic, 1, false, body); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (a *orbAgent) subscribeGroupChannels(channels []string) []string {
	var successList []string
	for _, channelID := range channels {

		base := fmt.Sprintf("channels/%s/messages", channelID)
		rpcFromCoreTopic := fmt.Sprintf("%s/%s", base, fleet.RPCFromCoreTopic)

		token := a.client.Subscribe(rpcFromCoreTopic, 1, a.handleGroupRPCFromCore)
		if token.Error() != nil {
			a.logger.Error("failed to subscribe to group channel/topic", zap.String("topic", rpcFromCoreTopic), zap.Error(token.Error()))
			continue
		}
		ok := token.WaitTimeout(time.Second * 5)
		if ok && token.Error() != nil {
			a.logger.Error("failed to subscribe to group channel/topic", zap.String("topic", rpcFromCoreTopic), zap.Error(token.Error()))
			continue
		}
		if !ok {
			a.logger.Error("failed to subscribe to group channel/topic: time out", zap.String("topic", rpcFromCoreTopic))
			continue
		}
		a.logger.Info("subscribed to a group", zap.String("topic", rpcFromCoreTopic))
		successList = append(successList, channelID)
	}
	return successList
}

func (a *orbAgent) handleGroupMembership(rpc fleet.GroupMembershipRPCPayload) {

	// if this is the full list, reset all group subscriptions and subscribed to this list
	if rpc.FullList {
		a.unsubscribeGroupChannels()
		a.groupChannels = a.subscribeGroupChannels(rpc.ChannelIDS)
	} else {
		// otherwise, just add these subscriptions to the existing list
		successList := a.subscribeGroupChannels(rpc.ChannelIDS)
		a.groupChannels = append(a.groupChannels, successList...)
	}
}

func (a *orbAgent) handleAgentPolicies(rpc []fleet.AgentPolicyRPCPayload) {

	a.logger.Info("Agent Policy from core", zap.Any("payload", rpc))

}

func (a *orbAgent) handleGroupRPCFromCore(client mqtt.Client, message mqtt.Message) {

	a.logger.Info("Group RPC message from core", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))

	var rpc fleet.RPC
	if err := json.Unmarshal(message.Payload(), &rpc); err != nil {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaMalformed))
		return
	}
	if rpc.SchemaVersion != fleet.CurrentRPCSchemaVersion {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaVersion))
		return
	}
	if rpc.Func == "" || rpc.Payload == nil {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaMalformed))
		return
	}

	// dispatch
	switch rpc.Func {
	case fleet.AgentPolicyRPCFunc:
		var r fleet.AgentPolicyRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding agent policy message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleAgentPolicies(r.Payload)
	default:
		a.logger.Warn("unsupported/unhandled core RPC, ignoring",
			zap.String("func", rpc.Func),
			zap.Any("payload", rpc.Payload))
	}

}

func (a *orbAgent) handleRPCFromCore(client mqtt.Client, message mqtt.Message) {

	a.logger.Info("RPC message from core", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))

	var rpc fleet.RPC
	if err := json.Unmarshal(message.Payload(), &rpc); err != nil {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaMalformed))
		return
	}
	if rpc.SchemaVersion != fleet.CurrentRPCSchemaVersion {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaVersion))
		return
	}
	if rpc.Func == "" || rpc.Payload == nil {
		a.logger.Error("error decoding RPC message from core", zap.Error(fleet.ErrSchemaMalformed))
		return
	}

	// dispatch
	switch rpc.Func {
	case fleet.GroupMembershipRPCFunc:
		var r fleet.GroupMembershipRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding group membership message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleGroupMembership(r.Payload)
	case fleet.AgentPolicyRPCFunc:
		var r fleet.AgentPolicyRPC
		if err := json.Unmarshal(message.Payload(), &r); err != nil {
			a.logger.Error("error decoding agent policy message from core", zap.Error(fleet.ErrSchemaMalformed))
			return
		}
		a.handleAgentPolicies(r.Payload)
	default:
		a.logger.Warn("unsupported/unhandled core RPC, ignoring",
			zap.String("func", rpc.Func),
			zap.Any("payload", rpc.Payload))
	}

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
