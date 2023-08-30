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
	opentracing "github.com/opentracing/opentracing-go"
	fleetgrpc "github.com/orb-community/orb/fleet/api/grpc"
	fleetpb "github.com/orb-community/orb/fleet/pb"
	"github.com/orb-community/orb/pkg/config"
	"github.com/orb-community/orb/policies"
	policiesgrpc "github.com/orb-community/orb/policies/api/grpc"
	policieshttp "github.com/orb-community/orb/policies/api/http"
	"github.com/orb-community/orb/policies/pb"
	"github.com/orb-community/orb/policies/postgres"
	rediscon "github.com/orb-community/orb/policies/redis/consumer"
	redisprod "github.com/orb-community/orb/policies/redis/producer"
	sinkgrpc "github.com/orb-community/orb/sinks/api/grpc"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/reflection"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	r "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/mainflux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	svcName     = "policies"
	mfEnvPrefix = "mf"
	envPrefix   = "orb_policies"
	httpPort    = "8202"
)

func main() {

	authGRPCCfg := config.LoadGRPCConfig(mfEnvPrefix, "auth")

	esCfg := config.LoadEsConfig(envPrefix)
	svcCfg := config.LoadBaseServiceConfig(envPrefix, httpPort)
	dbCfg := config.LoadPostgresConfig(envPrefix, svcName)
	jCfg := config.LoadJaegerConfig(envPrefix)
	fleetGRPCCfg := config.LoadGRPCConfig("orb", "fleet")
	policiesGRPCCfg := config.LoadGRPCConfig("orb", "policies")
	sinksGRPCCfg := config.LoadGRPCConfig("orb", "sinks")

	// logger
	var logger *zap.Logger
	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToLower(svcCfg.LogLevel) {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		os.Stdout,
		atomicLevel,
	)
	logger = zap.New(core, zap.AddCaller())
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	db := connectToDB(dbCfg, logger)
	defer db.Close()

	esClient := connectToRedis(esCfg.URL, esCfg.Pass, esCfg.DB, logger)
	defer esClient.Close()

	tracer, tracerCloser := initJaeger(svcName, jCfg.URL, logger)
	defer tracerCloser.Close()

	authGRPCConn := connectToGRPC(authGRPCCfg, logger)
	defer authGRPCConn.Close()

	authGRPCTimeout, err := time.ParseDuration(authGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", authGRPCCfg.Timeout, err.Error())
	}
	authGRPCClient := authapi.NewClient(tracer, authGRPCConn, authGRPCTimeout)

	fleetGRPCConn := connectToGRPC(fleetGRPCCfg, logger)
	defer fleetGRPCConn.Close()

	fleetGRPCTimeout, err := time.ParseDuration(fleetGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", fleetGRPCCfg.Timeout, err.Error())
	}

	fleetGRPCClient := fleetgrpc.NewClient(tracer, fleetGRPCConn, fleetGRPCTimeout)

	sinksGRPCConn := connectToGRPC(sinksGRPCCfg, logger)
	defer fleetGRPCConn.Close()

	sinksGRPCTimeout, err := time.ParseDuration(sinksGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", sinksGRPCCfg.Timeout, err.Error())
	}

	sinksGRPCClient := sinkgrpc.NewClient(tracer, sinksGRPCConn, sinksGRPCTimeout, logger)

	svc := newService(authGRPCClient, db, logger, esClient, fleetGRPCClient, sinksGRPCClient)
	errs := make(chan error, 2)

	go startHTTPServer(tracer, svc, svcCfg, logger, errs)
	go startGRPCServer(svc, tracer, policiesGRPCCfg, logger, errs)
	go subscribeToFleetES(svc, esClient, esCfg, logger)
	go subscribeToSinksES(svc, esClient, esCfg, logger)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Policies service terminated: %s", err))
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
		logger.Error("Failed to init Jaeger client, tracing will be disabled", zap.Error(err))
	}

	return tracer, closer
}

func newService(auth mainflux.AuthServiceClient, db *sqlx.DB, logger *zap.Logger, esClient *r.Client, fleetGrpcClient fleetpb.FleetServiceClient, sinksGrpcClient sinkspb.SinkServiceClient) policies.Service {
	thingsRepo := postgres.NewPoliciesRepository(db, logger)

	svc := policies.New(logger, auth, thingsRepo, fleetGrpcClient, sinksGrpcClient)
	svc = redisprod.NewEventStoreMiddleware(svc, esClient, logger)
	svc = policieshttp.NewLoggingMiddleware(svc, logger)
	svc = policieshttp.MetricsMiddleware(
		auth,
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "policies",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method", "owner_id", "policy_id", "dataset_id"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "policies",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method", "owner_id", "policy_id", "dataset_id"}),
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

func startHTTPServer(tracer opentracing.Tracer, svc policies.Service, cfg config.BaseSvcConfig, logger *zap.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.HttpPort)
	if cfg.HttpServerCert != "" || cfg.HttpServerKey != "" {
		logger.Info(fmt.Sprintf("Policies service started using https on port %s with cert %s key %s",
			cfg.HttpPort, cfg.HttpServerCert, cfg.HttpServerKey))
		errs <- http.ListenAndServeTLS(p, cfg.HttpServerCert, cfg.HttpServerKey, policieshttp.MakeHandler(tracer, svcName, svc))
		return
	}
	logger.Info(fmt.Sprintf("Policies service started using http on port %s", cfg.HttpPort))
	errs <- http.ListenAndServe(p, policieshttp.MakeHandler(tracer, svcName, svc))
}

func startGRPCServer(svc policies.Service, tracer opentracing.Tracer, cfg config.GRPCConfig, logger *zap.Logger, errs chan error) {
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

	pb.RegisterPolicyServiceServer(server, policiesgrpc.NewServer(tracer, svc))
	reflection.Register(server)
	errs <- server.Serve(listener)
}

func subscribeToFleetES(svc policies.Service, client *r.Client, cfg config.EsConfig, logger *zap.Logger) {
	eventStore := rediscon.NewEventStore(svc, client, cfg.Consumer, logger)
	logger.Info("Subscribed to Redis Event Store for agent groups")
	if err := eventStore.SubscribeToFleet(context.Background()); err != nil {
		logger.Error("Bootstrap service failed to subscribe to event sourcing", zap.Error(err))
	}
}

func subscribeToSinksES(svc policies.Service, client *r.Client, cfg config.EsConfig, logger *zap.Logger) {
	eventStore := rediscon.NewEventStore(svc, client, cfg.Consumer, logger)
	logger.Info("Subscribed to Redis Event Store for sinks")
	if err := eventStore.SubscribeToSink(context.Background()); err != nil {
		logger.Error("Bootstrap service failed to subscribe to event sourcing", zap.Error(err))
	}
}
