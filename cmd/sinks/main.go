// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	authapi "github.com/mainflux/mainflux/auth/api/grpc"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/pkg/config"
	"github.com/ns1labs/orb/sinks"
	sinksgrpc "github.com/ns1labs/orb/sinks/api/grpc"
	sinkshttp "github.com/ns1labs/orb/sinks/api/http"
	"github.com/ns1labs/orb/sinks/pb"
	"github.com/ns1labs/orb/sinks/postgres"
	redisprod "github.com/ns1labs/orb/sinks/redis/producer"
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
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	svcName     = "sinks"
	mfEnvPrefix = "mf"
	envPrefix   = "orb_sinks"
	httpPort    = "8200"
)

func main() {

	authCfg := config.LoadGRPCConfig(mfEnvPrefix, "auth")
	sdkCfg := config.LoadMFSDKConfig(mfEnvPrefix)

	esCfg := config.LoadEsConfig(envPrefix)
	svcCfg := config.LoadBaseServiceConfig(envPrefix, httpPort)
	dbCfg := config.LoadPostgresConfig(envPrefix, svcName)
	jCfg := config.LoadJaegerConfig(envPrefix)
	sinksGRPCCfg := config.LoadGRPCConfig("orb", "sinks")

	// main logger
	var logger *zap.Logger
	if svcCfg.LogLevel == "debug" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync() // flushes buffer, if any

	db := connectToDB(dbCfg, logger)
	defer db.Close()

	esClient := connectToRedis(esCfg.URL, esCfg.Pass, esCfg.DB, logger)
	defer esClient.Close()

	tracer, tracerCloser := initJaeger(svcName, jCfg.URL, logger)
	defer tracerCloser.Close()

	authConn := connectToAuth(authCfg, logger)
	defer authConn.Close()

	authTimeout, err := time.ParseDuration(authCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", authCfg.Timeout, err.Error())
	}
	auth := authapi.NewClient(tracer, authConn, authTimeout)

	sinkRepo := postgres.NewSinksRepository(db, logger)

	svc := newSinkService(auth, logger, esClient, sdkCfg, sinkRepo)
	errs := make(chan error, 2)

	go startHTTPServer(tracer, svc, svcCfg, logger, errs)
	go startGRPCServer(svc, tracer, sinksGRPCCfg, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Sink service terminated: %s", err))
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

func newSinkService(auth mainflux.AuthServiceClient, logger *zap.Logger, esClient *r.Client, sdkCfg config.MFSDKConfig, repoSink sinks.SinkRepository) sinks.SinkService {

	config := mfsdk.Config{
		BaseURL:      sdkCfg.BaseURL,
		ThingsPrefix: sdkCfg.ThingsPrefix,
	}

	mfsdk := mfsdk.NewSDK(config)

	svc := sinks.NewSinkService(logger, auth, repoSink, mfsdk)
	svc = redisprod.NewEventStoreMiddleware(svc, esClient)
	svc = sinkshttp.NewLoggingMiddleware(svc, logger)
	svc = sinkshttp.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "sink",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "sink",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func connectToAuth(cfg config.GRPCConfig, logger *zap.Logger) *grpc.ClientConn {
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
		logger.Info("gRPC communication is not encrypted")
	}

	conn, err := grpc.Dial(cfg.URL, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to auth service: %s", err))
		os.Exit(1)
	}

	return conn
}

func startHTTPServer(tracer opentracing.Tracer, svc sinks.SinkService, cfg config.BaseSvcConfig, logger *zap.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", cfg.HttpPort)
	if cfg.HttpServerCert != "" || cfg.HttpServerKey != "" {
		logger.Info(fmt.Sprintf("Sink service started using https on port %s with cert %s key %s",
			cfg.HttpPort, cfg.HttpServerCert, cfg.HttpServerKey))
		errs <- http.ListenAndServeTLS(p, cfg.HttpServerCert, cfg.HttpServerKey, sinkshttp.MakeHandler(tracer, svcName, svc))
		return
	}
	logger.Info(fmt.Sprintf("Sink service started using http on port %s", cfg.HttpPort))
	errs <- http.ListenAndServe(p, sinkshttp.MakeHandler(tracer, svcName, svc))
}

func startGRPCServer(svc sinks.SinkService, tracer opentracing.Tracer, cfg config.GRPCConfig, logger *zap.Logger, errs chan error) {
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
	pb.RegisterSinkServiceServer(server, sinksgrpc.NewServer(tracer, svc))
	reflection.Register(server)
	errs <- server.Serve(listener)
}
