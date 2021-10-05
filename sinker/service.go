/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mainflux/mainflux/pkg/messaging"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/ns1labs/orb/fleet"
	fleetpb "github.com/ns1labs/orb/fleet/pb"
	policiespb "github.com/ns1labs/orb/policies/pb"
	"github.com/ns1labs/orb/sinker/backend"
	"github.com/ns1labs/orb/sinker/backend/pktvisor"
	"github.com/ns1labs/orb/sinker/config"
	"github.com/ns1labs/orb/sinker/prometheus"
	sinkspb "github.com/ns1labs/orb/sinks/pb"
	"go.uber.org/zap"
	"os"
	"strings"
)

const (
	BackendMetricsTopic = "be.*.m.>"
	MaxMsgPayloadSize   = 1024 * 100
)

var (
	ErrPayloadTooBig = errors.New("payload too big")
)

type Service interface {
	// Start set up communication with the message bus to communicate with agents
	Start() error
	// Stop end communication with the message bus
	Stop() error
}

type sinkerService struct {
	pubSub mfnats.PubSub

	esclient *redis.Client

	logger *zap.Logger

	configRepo config.ConfigRepo

	promClient prometheus.Client

	policiesClient policiespb.PolicyServiceClient
	fleetClient    fleetpb.FleetServiceClient
	sinksClient    sinkspb.SinkServiceClient
}

func (svc sinkerService) remoteWriteToPrometheus(tsList prometheus.TSList) {

	svc.logger.Info("writing to", zap.String("url", prometheus.DefaultRemoteWrite))

	result, writeErr := svc.promClient.WriteTimeSeries(context.Background(), tsList,
		prometheus.WriteOptions{})
	if err := error(writeErr); err != nil {
		json.NewEncoder(os.Stdout).Encode(struct {
			Success    bool   `json:"success"`
			Error      string `json:"error"`
			StatusCode int    `json:"statusCode"`
		}{
			Success:    false,
			Error:      err.Error(),
			StatusCode: writeErr.StatusCode(),
		})
		os.Stdout.Sync()

		svc.logger.Error("remote write error", zap.Error(err))
	}

	json.NewEncoder(os.Stdout).Encode(struct {
		Success    bool `json:"success"`
		StatusCode int  `json:"statusCode"`
	}{
		Success:    true,
		StatusCode: result.StatusCode,
	})
	os.Stdout.Sync()

	svc.logger.Info("write success")
}

func (svc sinkerService) handleSinkConfig(channelID string, metrics []fleet.AgentMetricsRPCPayload) error {
	//Todo use channelID to get owner id on fleet grpc server
	ownerID, err := svc.fleetClient.RetrieveOwnerByChannelID(context.Background(), &fleetpb.OwnerByChannelIDReq{Channel: channelID})
	if err != nil {
		return err
	}
	fmt.Printf("Got this owner %s by this channel: %s", ownerID, channelID)
	//Todo use metricsRPC.Payload[0].Datasets[] to retrieve the sinkID from policy grpc server
	//var mapDataset = map[string][]string{}
	for _, m := range metrics {
		//var sinkIDs []string
		for _, ds := range m.Datasets {
			fmt.Sprintf(ds)
			//sinkID, err := svc.policiesClient.
		}
	}

	//Todo use the retrieve sinkID to get the backend config

	return nil
}

func (svc sinkerService) handleMetrics(thingID string, channelID string, subtopic string, payload []byte) error {
	// find backend to send it to
	s := strings.Split(subtopic, ".")
	if len(s) < 3 || s[0] != "be" || s[2] != "m" {
		return errors.New(fmt.Sprintf("invalid subtopic, ignoring: %s", subtopic))
	}
	if !backend.HaveBackend(s[1]) {
		return errors.New(fmt.Sprintf("unknown agent backend, ignoring: %s", s[1]))
	}
	be := backend.GetBackend(s[1])
	// unpack metrics RPC
	var versionCheck fleet.SchemaVersionCheck
	if err := json.Unmarshal(payload, &versionCheck); err != nil {
		return fleet.ErrSchemaMalformed
	}
	if versionCheck.SchemaVersion != fleet.CurrentRPCSchemaVersion {
		return fleet.ErrSchemaVersion
	}
	var rpc fleet.RPC
	if err := json.Unmarshal(payload, &rpc); err != nil {
		return fleet.ErrSchemaMalformed
	}
	if rpc.Func != fleet.AgentMetricsRPCFunc {
		return errors.New(fmt.Sprintf("unexpected RPC function: %s", rpc.Func))
	}
	var metricsRPC fleet.AgentMetricsRPC
	if err := json.Unmarshal(payload, &metricsRPC); err != nil {
		return fleet.ErrSchemaMalformed
	}

	if err := svc.handleSinkConfig(channelID, metricsRPC.Payload); err != nil {
		return err
	}

	tsList, err := be.ProcessMetrics(thingID, channelID, s, metricsRPC.Payload)
	if err != nil {
		return err
	}

	svc.remoteWriteToPrometheus(tsList)

	return nil
}

func (svc sinkerService) handleMsgFromAgent(msg messaging.Message) error {

	// NOTE: we need to consider ALL input from the agent as untrusted, the same as untrusted HTTP API would be

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}

	svc.logger.Debug("received agent message",
		zap.Any("payload", payload),
		zap.String("subtopic", msg.Subtopic),
		zap.String("channel", msg.Channel),
		zap.String("protocol", msg.Protocol),
		zap.Int64("created", msg.Created),
		zap.String("publisher", msg.Publisher))

	if len(msg.Payload) > MaxMsgPayloadSize {
		return ErrPayloadTooBig
	}

	if err := svc.handleMetrics(msg.Publisher, msg.Channel, msg.Subtopic, msg.Payload); err != nil {
		svc.logger.Error("metrics processing failure", zap.Error(err))
		return err
	}

	return nil
}

func (svc sinkerService) Start() error {

	topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
	if err := svc.pubSub.Subscribe(topic, svc.handleMsgFromAgent); err != nil {
		return err
	}
	svc.logger.Info("started metrics consumer", zap.String("topic", topic))
	return nil
}

func (svc sinkerService) Stop() error {
	topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
	if err := svc.pubSub.Unsubscribe(topic); err != nil {
		return err
	}
	svc.logger.Info("unsubscribed from agent metrics")
	return nil
}

// New instantiates the sinker service implementation.
func New(logger *zap.Logger,
	pubSub mfnats.PubSub,
	esclient *redis.Client,
	configRepo config.ConfigRepo,
	policiesClient policiespb.PolicyServiceClient,
	fleetClient fleetpb.FleetServiceClient,
	sinksClient sinkspb.SinkServiceClient,
	promClient prometheus.Client) Service {

	pktvisor.Register(logger)
	return &sinkerService{
		logger:         logger,
		pubSub:         pubSub,
		esclient:       esclient,
		configRepo:     configRepo,
		policiesClient: policiesClient,
		fleetClient:    fleetClient,
		sinksClient:    sinksClient,
		promClient:     promClient,
	}
}
