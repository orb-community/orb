// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package maestro

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/maestro/kubecontrol"
	rediscons1 "github.com/ns1labs/orb/maestro/redis/consumer"
	"github.com/ns1labs/orb/pkg/config"
	sinkspb "github.com/ns1labs/orb/sinks/pb"
	"go.uber.org/zap"
	"time"
)

var _ Service = (*maestroService)(nil)

type maestroService struct {
	serviceContext    context.Context
	serviceCancelFunc context.CancelFunc

	kubecontrol kubecontrol.Service
	logger      *zap.Logger
	redisClient *redis.Client
	sinksClient sinkspb.SinkServiceClient
	esCfg       config.EsConfig
	eventStore  rediscons1.Subscriber
}

func NewMaestroService(logger *zap.Logger, redisClient *redis.Client, sinksGrpcClient sinkspb.SinkServiceClient, esCfg config.EsConfig) Service {
	kubectr := kubecontrol.NewService(logger)
	eventStore := rediscons1.NewEventStore(redisClient, kubectr, esCfg.Consumer, logger)
	return &maestroService{
		logger:      logger,
		redisClient: redisClient,
		sinksClient: sinksGrpcClient,
		kubecontrol: kubectr,
		eventStore:  eventStore,
	}
}

type SinkData struct {
	SinkID          string          `json:"sink_id"`
	OwnerID         string          `json:"owner_id"`
	Url             string          `json:"remote_host"`
	User            string          `json:"username"`
	Password        string          `json:"password"`
	State           PrometheusState `json:"state,omitempty"`
	Msg             string          `json:"msg,omitempty"`
	LastRemoteWrite time.Time       `json:"last_remote_write,omitempty"`
}

const (
	Unknown PrometheusState = iota
	Active
	Error
	Idle
)

type PrometheusState int

var promStateMap = [...]string{
	"unknown",
	"active",
	"error",
	"idle",
}

var promStateRevMap = map[string]PrometheusState{
	"unknown": Unknown,
	"active":  Active,
	"error":   Error,
	"idle":    Idle,
}

func (p PrometheusState) String() string {
	return promStateMap[p]
}

func (p *PrometheusState) Scan(value interface{}) error {
	*p = promStateRevMap[string(value.([]byte))]
	return nil
}

func (p PrometheusState) Value() (driver.Value, error) {
	return p.String(), nil
}

// Start will load all sinks from DB using SinksGRPC,
//
//	then for each sink, will create DeploymentEntry in Redis
//	And for each sink with active state, deploy OtelCollector
func (svc *maestroService) Start(ctx context.Context, cancelFunction context.CancelFunc) error {

	loadCtx, loadCancelFunction := context.WithCancel(ctx)
	defer loadCancelFunction()
	svc.serviceContext = ctx
	svc.serviceCancelFunc = cancelFunction

	sinksRes, err := svc.sinksClient.RetrieveSinks(loadCtx, &sinkspb.SinksFilterReq{OtelEnabled: "enabled"})
	if err != nil {
		loadCancelFunction()
		return err
	}

	for _, sinkRes := range sinksRes.Sinks {
		sinkContext := context.WithValue(loadCtx, "sink-id", sinkRes.Id)
		var data SinkData
		if err := json.Unmarshal(sinkRes.Config, &data); err != nil {
			svc.logger.Warn("failed to unmarshal sink, skipping", zap.String("sink-id", data.SinkID))
			continue
		}

		err := svc.eventStore.CreateDeploymentEntry(sinkContext, sinkRes.Id, data.Url, data.User, data.Password)
		if err != nil {
			svc.logger.Warn("failed to create deploymentEntry for sink, skipping", zap.String("sink-id", data.SinkID))
			continue
		}

		// if State is Active, deploy OtelCollector
		if data.State == Active {
			deploymentEntry, err := svc.eventStore.GetDeploymentEntryFromSinkId(sinkContext, data.SinkID)
			if err != nil {
				svc.logger.Warn("failed to fetch deploymentEntry for sink, skipping", zap.String("sink-id", data.SinkID))
				continue
			}
			err = svc.kubecontrol.CreateOtelCollector(sinkContext, data.SinkID, deploymentEntry)
			if err != nil {
				svc.logger.Warn("failed to deploy OtelCollector for sink, skipping", zap.String("sink-id", data.SinkID))
				continue
			}
		}
	}

	go svc.subscribeToSinksES()
	go svc.subscribeToSinkerES()
	return nil
}

func (svc *maestroService) subscribeToSinkerES() {
	if err := svc.eventStore.SubscribeSinker(context.Background()); err != nil {
		svc.logger.Error("Bootstrap service failed to subscribe to event sourcing sinker", zap.Error(err))
	}
	svc.logger.Info("Subscribed to Redis Event Store for sinker")
}

func (svc *maestroService) subscribeToSinksES() {
	svc.logger.Info("Subscribed to Redis Event Store for sinks")
	if err := svc.eventStore.SubscribeSinks(context.Background()); err != nil {
		svc.logger.Error("Bootstrap service failed to subscribe to event sourcing sinks", zap.Error(err))
	}
}
