// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package maestro

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/orb-community/orb/maestro/deployment"
	"github.com/orb-community/orb/maestro/kubecontrol"
	"github.com/orb-community/orb/maestro/monitor"
	rediscons1 "github.com/orb-community/orb/maestro/redis/consumer"
	"github.com/orb-community/orb/maestro/redis/producer"
	"github.com/orb-community/orb/maestro/service"
	"github.com/orb-community/orb/pkg/config"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
)

var _ Service = (*maestroService)(nil)

type maestroService struct {
	serviceContext    context.Context
	serviceCancelFunc context.CancelFunc

	deploymentService   deployment.Service
	sinkListenerService rediscons1.SinksListener
	activityListener    rediscons1.SinkerActivityListener

	kubecontrol       kubecontrol.Service
	monitor           monitor.Service
	logger            *zap.Logger
	streamRedisClient *redis.Client
	sinkerRedisClient *redis.Client
	sinksClient       sinkspb.SinkServiceClient
	eventService      service.EventService
	esCfg             config.EsConfig
	kafkaUrl          string
}

func NewMaestroService(logger *zap.Logger, streamRedisClient *redis.Client, sinkerRedisClient *redis.Client,
	sinksGrpcClient sinkspb.SinkServiceClient, esCfg config.EsConfig, otelCfg config.OtelConfig, db *sqlx.DB) Service {
	kubectr := kubecontrol.NewService(logger)
	repo := deployment.NewRepositoryService(db, logger)
	deploymentService := deployment.NewDeploymentService(logger, repo, otelCfg.KafkaUrl, esCfg.EncryptionKey)
	ps := producer.NewMaestroProducer(logger, streamRedisClient)
	monitorService := monitor.NewMonitorService(logger, &sinksGrpcClient, ps, &kubectr)
	eventService := service.NewEventService(logger, deploymentService, kubectr)
	sinkListenerService := rediscons1.NewSinksListenerController(logger, eventService, sinkerRedisClient, sinksGrpcClient)
	activityListener := rediscons1.NewSinkerActivityListener(logger, eventService, sinkerRedisClient)
	return &maestroService{
		logger:              logger,
		deploymentService:   deploymentService,
		streamRedisClient:   streamRedisClient,
		sinkerRedisClient:   sinkerRedisClient,
		sinksClient:         sinksGrpcClient,
		sinkListenerService: sinkListenerService,
		activityListener:    activityListener,
		kubecontrol:         kubectr,
		monitor:             monitorService,
		kafkaUrl:            otelCfg.KafkaUrl,
	}
}

// Start will load all sinks from DB using SinksGRPC,
//
//	then for each sink, will create DeploymentEntry in Redis
//	And for each sink with active state, deploy OtelCollector
func (svc *maestroService) Start(ctx context.Context, cancelFunction context.CancelFunc) error {

	svc.serviceContext = ctx
	svc.serviceCancelFunc = cancelFunction

	go svc.subscribeToSinksEvents(ctx)
	go svc.subscribeToSinkerEvents(ctx)

	monitorCtx := context.WithValue(ctx, "routine", "monitor")
	err := svc.monitor.Start(monitorCtx, cancelFunction)
	if err != nil {
		svc.logger.Error("error during monitor routine start", zap.Error(err))
		cancelFunction()
		return err
	}
	svc.logger.Info("Maestro service started")

	return nil
}

func (svc *maestroService) Stop() {
	svc.serviceCancelFunc()
	svc.logger.Info("Maestro service stopped")
}

func (svc *maestroService) subscribeToSinksEvents(ctx context.Context) {
	if err := svc.sinkListenerService.SubscribeSinksEvents(ctx); err != nil {
		svc.logger.Error("Bootstrap service failed to subscribe to event sourcing", zap.Error(err))
		return
	}
	svc.logger.Info("finished reading sinks events")
	ctx.Done()
}

func (svc *maestroService) subscribeToSinkerEvents(ctx context.Context) {
	if err := svc.activityListener.SubscribeSinksEvents(ctx); err != nil {
		svc.logger.Error("Bootstrap service failed to subscribe to event sourcing", zap.Error(err))
		return
	}
	svc.logger.Info("finished reading sinker events")
	ctx.Done()
}
