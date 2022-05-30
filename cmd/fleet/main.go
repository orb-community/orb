// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"context"
	"fmt"
	authapi "github.com/mainflux/mainflux/auth/api/grpc"
	mflog "github.com/mainflux/mainflux/logger"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/ns1labs/orb/fleet"
	fleetgrpc "github.com/ns1labs/orb/fleet/api/grpc"
	fleethttp "github.com/ns1labs/orb/fleet/api/http"
	"github.com/ns1labs/orb/fleet/backend/pktvisor"
	"github.com/ns1labs/orb/fleet/pb"
	"github.com/ns1labs/orb/fleet/postgres"
	rediscons "github.com/ns1labs/orb/fleet/redis/consumer"
	redisprod "github.com/ns1labs/orb/fleet/redis/producer"
	"github.com/ns1labs/orb/pkg/config"
	policiesgrpc "github.com/ns1labs/orb/policies/api/grpc"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	r "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	svcName     = "fleet"
	mfEnvPrefix = "mf"
	envPrefix   = "orb_fleet"
	httpPort    = "8203"
)

func main() {

	natsCfg := config.LoadNatsConfig(envPrefix)
	authGRPCCfg := config.LoadGRPCConfig(mfEnvPrefix, "auth")
	sdkCfg := config.LoadMFSDKConfig(mfEnvPrefix)

	esCfg := config.LoadEsConfig(envPrefix)
	svcCfg := config.LoadBaseServiceConfig(envPrefix, httpPort)
	dbCfg := config.LoadPostgresConfig(envPrefix, svcName)
	jCfg := config.LoadJaegerConfig(envPrefix)
	policiesGRPCCfg := config.LoadGRPCConfig("orb", "policies")
	fleetGRPCCfg := config.LoadGRPCConfig("orb", "fleet")

	// main logger
	var logger *zap.Logger
	if svcCfg.LogLevel == "debug" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync() // flushes buffer, if any

	// only needed for mainflux interfaces
	mflogger, err := mflog.New(os.Stdout, svcCfg.LogLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := connectToDB(dbCfg, logger)
	defer db.Close()

	esClient := connectToRedis(esCfg.URL, esCfg.Pass, esCfg.DB, logger)
	defer esClient.Close()

	tracer, tracerCloser := initJaeger(svcName, jCfg.URL, logger)
	defer tracerCloser.Close()

	authGRPCConn := connectToGRPC(authGRPCCfg, logger)
	defer authGRPCConn.Close()

	policiesGRPCConn := connectToGRPC(policiesGRPCCfg, logger)
	defer policiesGRPCConn.Close()

	authGRPCTimeout, err := time.ParseDuration(authGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", authGRPCCfg.Timeout, err.Error())
	}
	authGRPCClient := authapi.NewClient(tracer, authGRPCConn, authGRPCTimeout)

	policiesGRPCTimeout, err := time.ParseDuration(policiesGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", policiesGRPCCfg.Timeout, err.Error())
	}
	policiesGRPCClient := policiesgrpc.NewClient(tracer, policiesGRPCConn, policiesGRPCTimeout)

	pubSub, err := mfnats.NewPubSub(natsCfg.URL, svcName, mflogger)
	if err != nil {
		logger.Error("Failed to connect to NATS", zap.Error(err))
		os.Exit(1)
	}
	defer pubSub.Close()

	agentRepo := postgres.NewAgentRepository(db, logger)
	agentGroupRepo := postgres.NewAgentGroupRepository(db, logger)

	commsSvc := fleet.NewFleetCommsService(logger, policiesGRPCClient, agentRepo, agentGroupRepo, pubSub)
	commsSvc = fleet.CommsMetricsMiddleware(
		commsSvc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "fleet",
			Subsystem: "comms",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method", "agent_id", "agent_name", "group_id", "group_name", "owner_id"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "fleet",
			Subsystem: "comms",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method", "agent_id", "agent_name", "group_id", "group_name", "owner_id"}))

	aDone := make(chan bool)

	svc := newFleetService(authGRPCClient, db, logger, esClient, sdkCfg, agentRepo, agentGroupRepo, commsSvc, aDone)
	defer commsSvc.Stop()

	errs := make(chan error, 2)

	go startHTTPServer(tracer, svc, svcCfg, logger, errs)
	go subscribeToPoliciesES(svc, commsSvc, esClient, esCfg, logger)
	go startGRPCServer(svc, tracer, fleetGRPCCfg, logger, errs)

	err = commsSvc.Start()
	if err != nil {
		logger.Error("unable to start agent communication", zap.Error(err))
		os.Exit(1)
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
		aDone <- true
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Fleet service terminated: %s", err))
}

func connectToDB(cfg config.PostgresConfig, logger *zap.Logger) *sqlx.DB {
	db, err := postgres.Connect(cfg)
	if err != nil {
		logger.Error("Failed to connect to postgres", zap.Error(err))
		os.Exit(1)
	}
	return db
}

func connectToRedis(redisURL, redisPass, redisDB string, logger *zap.Logger) *r.Client {
	db, err := strconv.Atoi(redisDB)
	if err != nil {
		logger.Error("Failed to connect to redis", zap.Error(err))
		os.Exit(1)
	}

	return r.NewClient(&r.Options{
		Addr:     redisURL,
		Password: redisPass,
		DB:       db,
	})
}

func initJaeger(svcName, url string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
	if url == "" {
		return opentracing.NoopTracer{}, ioutil.NopCloser(nil)
	}

	tracer, closer, err := jconfig.Configuration{
		ServiceName: svcName,
		Sampler: &jconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: url,
			LogSpans:           true,
		},
	}.NewTracer()
	if err != nil {
		logger.Error("Failed to init Jaeger client", zap.Error(err))
		os.Exit(1)
	}

	return tracer, closer
}

func newFleetService(auth mainflux.AuthServiceClient, db *sqlx.DB, logger *zap.Logger, esClient *r.Client, sdkCfg config.MFSDKConfig, agentRepo fleet.AgentRepository, agentGroupRepo fleet.AgentGroupRepository, agentComms fleet.AgentCommsService, aDone chan bool) fleet.Service {

	config := mfsdk.Config{
		ThingsURL: fmt.Sprintf("%s%s", sdkCfg.BaseURL, sdkCfg.ThingsPrefix),
	}

	mfsdk := mfsdk.NewSDK(config)

	pktvisor.Register(auth, agentRepo)

	svc := fleet.NewFleetService(logger, auth, agentRepo, agentGroupRepo, agentComms, mfsdk, aDone)
	svc = redisprod.NewEventStoreMiddleware(svc, esClient)
	svc = fleethttp.NewLoggingMiddleware(svc, logger)
	svc = fleethttp.MetricsMiddleware(
		auth,
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "fleet",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method", "owner_id", "agent_id", "group_id"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "fleet",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method", "owner_id", "agent_id", "group_id"}),
	)
	return svc
}

func connectToGRPC(cfg config.GRPCConfig, logger *zap.Logger) *grpc.ClientConn {
	var opts []grpc.DialOption
	tls, err := strconv.ParseBool(cfg.ClientTLS)
	if err != nil {
		tls = false
	}
	if tls {
		if cfg.CaCerts != "" {
			tpc, err := credentials.NewClientTLSFromFile(cfg.CaCerts, "")
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to create tls credentials: %s", err))
				os.Exit(1)
			}
			opts = append(opts, grpc.WithTransportCredentials(tpc))
		}
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(cfg.URL, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to dial to gRPC service %s: %s", cfg.URL, err))
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Dialed to gRPC service %s at %s, TLS? %t", cfg.Service, cfg.URL, tls))

	return conn
}

func startHTTPServer(tracer opentracing.Tracer, svc fleet.Service, cfg config.BaseSvcConfig, logger *zap.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.HttpPort)
	if cfg.HttpServerCert != "" || cfg.HttpServerKey != "" {
		logger.Info(fmt.Sprintf("Fleet service started using https on port %s with cert %s key %s",
			cfg.HttpPort, cfg.HttpServerCert, cfg.HttpServerKey))
		errs <- http.ListenAndServeTLS(p, cfg.HttpServerCert, cfg.HttpServerKey, fleethttp.MakeHandler(tracer, svcName, svc))
		return
	}
	logger.Info(fmt.Sprintf("Fleet service started using http on port %s", cfg.HttpPort))
	errs <- http.ListenAndServe(p, fleethttp.MakeHandler(tracer, svcName, svc))
}

func subscribeToPoliciesES(svc fleet.Service, commsSvc fleet.AgentCommsService, client *r.Client, cfg config.EsConfig, logger *zap.Logger) {
	eventStore := rediscons.NewEventStore(svc, commsSvc, client, cfg.Consumer, logger)
	logger.Info("Subscribed to Redis Event Store for policies")
	if err := eventStore.Subscribe(context.Background()); err != nil {
		logger.Error("Bootstrap service failed to subscribe to event sourcing", zap.Error(err))
	}
}

func startGRPCServer(svc fleet.Service, tracer opentracing.Tracer, cfg config.GRPCConfig, logger *zap.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.Port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to gRPC listen on port %s: %s", cfg.Port, err))
		os.Exit(1)
	}

	var server *grpc.Server
	if cfg.ServerCert != "" || cfg.ServerKey != "" {
		creds, err := credentials.NewServerTLSFromFile(cfg.ServerCert, cfg.ServerKey)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to load things certificates: %s", err))
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("gRPC service started using https on port %s with cert %s key %s",
			cfg.Port, cfg.ServerCert, cfg.ServerKey))
		server = grpc.NewServer(grpc.Creds(creds))
	} else {
		logger.Info(fmt.Sprintf("gRPC service started using http on port %s", cfg.Port))
		server = grpc.NewServer()
	}

	pb.RegisterFleetServiceServer(server, fleetgrpc.NewServer(tracer, svc))
	reflection.Register(server)
	errs <- server.Serve(listener)
}
