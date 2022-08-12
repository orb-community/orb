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
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-redis/redis/v8"
	"github.com/go-zoo/bone"
	"github.com/ns1labs/orb/buildinfo"
	fleetgrpc "github.com/ns1labs/orb/fleet/api/grpc"
	"github.com/ns1labs/orb/pkg/config"
	policiesgrpc "github.com/ns1labs/orb/policies/api/grpc"
	"github.com/ns1labs/orb/sinker"
	sinkconfig "github.com/ns1labs/orb/sinker/config"
	cacheconfig "github.com/ns1labs/orb/sinker/redis"
	"github.com/ns1labs/orb/sinker/redis/consumer"
	"github.com/ns1labs/orb/sinker/redis/producer"
	sinksgrpc "github.com/ns1labs/orb/sinks/api/grpc"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jconfig "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	mflog "github.com/mainflux/mainflux/logger"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
)

const (
	svcName   = "sinker"
	envPrefix = "orb_sinker"
	httpPort  = "8201"
)

func main() {

	natsCfg := config.LoadNatsConfig(envPrefix)
	esCfg := config.LoadEsConfig(envPrefix)
	cacheCfg := config.LoadCacheConfig(envPrefix)
	svcCfg := config.LoadBaseServiceConfig(envPrefix, httpPort)
	jCfg := config.LoadJaegerConfig(envPrefix)
	fleetGRPCCfg := config.LoadGRPCConfig("orb", "fleet")
	policiesGRPCCfg := config.LoadGRPCConfig("orb", "policies")
	sinksGRPCCfg := config.LoadGRPCConfig("orb", "sinks")

	// main logger
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
		err := logger.Sync()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(logger) // flushes buffer, if any

	// only needed for mainflux interfaces
	mflogger, err := mflog.New(os.Stdout, svcCfg.LogLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	cacheClient := connectToRedis(cacheCfg.URL, cacheCfg.Pass, cacheCfg.DB, logger)

	esClient := connectToRedis(esCfg.URL, esCfg.Pass, esCfg.DB, logger)
	defer func(esClient *redis.Client) {
		err := esClient.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(esClient)

	tracer, tracerCloser := initJaeger(svcName, jCfg.URL, logger)
	defer func(tracerCloser io.Closer) {
		err := tracerCloser.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(tracerCloser)

	pubSub, err := mfnats.NewPubSub(natsCfg.URL, svcName, mflogger)
	if err != nil {
		logger.Error("Failed to connect to NATS", zap.Error(err))
		os.Exit(1)
	}
	defer pubSub.Close()

	policiesGRPCConn := connectToGRPC(policiesGRPCCfg, logger)
	defer func(policiesGRPCConn *grpc.ClientConn) {
		err := policiesGRPCConn.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(policiesGRPCConn)

	policiesGRPCTimeout, err := time.ParseDuration(policiesGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", policiesGRPCCfg.Timeout, err.Error())
	}
	policiesGRPCClient := policiesgrpc.NewClient(tracer, policiesGRPCConn, policiesGRPCTimeout)

	fleetGRPCConn := connectToGRPC(fleetGRPCCfg, logger)
	defer func(fleetGRPCConn *grpc.ClientConn) {
		err := fleetGRPCConn.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(fleetGRPCConn)

	fleetGRPCTimeout, err := time.ParseDuration(fleetGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", fleetGRPCCfg.Timeout, err.Error())
	}
	fleetGRPCClient := fleetgrpc.NewClient(tracer, fleetGRPCConn, fleetGRPCTimeout)

	sinksGRPCConn := connectToGRPC(sinksGRPCCfg, logger)
	defer func(sinksGRPCConn *grpc.ClientConn) {
		err := sinksGRPCConn.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(sinksGRPCConn)

	sinksGRPCTimeout, err := time.ParseDuration(sinksGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", sinksGRPCCfg.Timeout, err.Error())
	}
	sinksGRPCClient := sinksgrpc.NewClient(tracer, sinksGRPCConn, sinksGRPCTimeout)

	configRepo := cacheconfig.NewSinkerCache(cacheClient, logger)
	configRepo = producer.NewEventStoreMiddleware(configRepo, esClient)
	gauge := kitprometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: "sinker",
		Subsystem: "sink",
		Name:      "payload_size",
		Help:      "Total size of outbound payloads",
	}, []string{"method", "agent_id", "agent", "policy_id", "policy", "sink_id", "owner_id"})
	counter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "sinker",
		Subsystem: "sink",
		Name:      "payload_count",
		Help:      "Number of payloads wrote",
	}, []string{"method", "agent_id", "agent", "policy_id", "policy", "sink_id", "owner_id"})
	inputCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "sinker",
		Subsystem: "sink",
		Name:      "message_inbound",
		Help:      "Number of messages received",
	}, []string{"subtopic", "channel", "protocol", "created", "publisher", "trace-id"})

	svc := sinker.New(logger, pubSub, esClient, configRepo, policiesGRPCClient, fleetGRPCClient, sinksGRPCClient, gauge, counter, inputCounter)
	defer func(svc sinker.Service) {
		err := svc.Stop()
		if err != nil {
			log.Fatalf("fatal error in stop the service: %e", err)
		}
	}(svc)

	errs := make(chan error, 2)

	go startHTTPServer(svcCfg, errs, logger)
	go subscribeToSinksES(svc, configRepo, esClient, esCfg, logger)

	err = svc.Start()
	if err != nil {
		logger.Error("unable to start agent metric consumption", zap.Error(err))
		os.Exit(1)
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error("sinker service terminated", zap.Error(err))
}

func makeHandler(svcName string) http.Handler {
	r := bone.New()
	r.GetFunc("/version", buildinfo.Version(svcName))
	r.Handle("/metrics", promhttp.Handler())
	return r
}

func startHTTPServer(cfg config.BaseSvcConfig, errs chan error, logger *zap.Logger) {
	p := fmt.Sprintf(":%s", cfg.HttpPort)
	if cfg.HttpServerCert != "" || cfg.HttpServerKey != "" {
		logger.Info(fmt.Sprintf("Sinker service started using https on port %s with cert %s key %s",
			cfg.HttpPort, cfg.HttpServerCert, cfg.HttpServerKey))
		errs <- http.ListenAndServeTLS(p, cfg.HttpServerCert, cfg.HttpServerKey, makeHandler(svcName))
		return
	}
	logger.Info(fmt.Sprintf("Sinker service started using http on port %s", cfg.HttpPort))
	errs <- http.ListenAndServe(p, makeHandler(svcName))
}

func connectToRedis(URL, pass string, cacheDB string, logger *zap.Logger) *redis.Client {
	db, err := strconv.Atoi(cacheDB)
	if err != nil {
		logger.Error("Failed to connect to cache", zap.Error(err))
		os.Exit(1)
	}

	return redis.NewClient(&redis.Options{
		Addr:     URL,
		Password: pass,
		DB:       db,
	})
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
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(cfg.URL, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to dial to gRPC service %s: %s", cfg.URL, err))
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Dialed to gRPC service %s at %s, TLS? %t", cfg.Service, cfg.URL, tls))

	return conn
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

func subscribeToSinksES(svc sinker.Service, configRepo sinkconfig.ConfigRepo, client *redis.Client, cfg config.EsConfig, logger *zap.Logger) {
	eventStore := consumer.NewEventStore(svc, configRepo, client, cfg.Consumer, logger)
	logger.Info("Subscribed to Redis Event Store for sinks")
	if err := eventStore.Subscribe(context.Background()); err != nil {
		logger.Error("Bootstrap service failed to subscribe to event sourcing", zap.Error(err))
	}
}
