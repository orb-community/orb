// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-zoo/bone"
	"github.com/ns1labs/orb"
	fleetgrpc "github.com/ns1labs/orb/fleet/api/grpc"
	"github.com/ns1labs/orb/pkg/config"
	policiesgrpc "github.com/ns1labs/orb/policies/api/grpc"
	"github.com/ns1labs/orb/sinker"
	config2 "github.com/ns1labs/orb/sinker/config"
	"github.com/ns1labs/orb/sinker/prometheus"
	sinksgrpc "github.com/ns1labs/orb/sinks/api/grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jconfig "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	var writeURLFlag string
	flag.StringVar(&writeURLFlag, "u", prometheus.DefaultRemoteWrite, "remote write endpoint")

	natsCfg := config.LoadNatsConfig(envPrefix)
	esCfg := config.LoadEsConfig(envPrefix)
	svcCfg := config.LoadBaseServiceConfig(envPrefix, httpPort)
	jCfg := config.LoadJaegerConfig(envPrefix)
	fleetGRPCCfg := config.LoadGRPCConfig("orb", "fleet")
	policiesGRPCCfg := config.LoadGRPCConfig("orb", "policies")
	sinksGRPCCfg := config.LoadGRPCConfig("orb", "sinks")

	cfg := prometheus.NewConfig(
		prometheus.WriteURLOption(writeURLFlag),
	)

	promClient, err := prometheus.NewClient(cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to construct client: %v", err))
	}

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

	esClient := connectToRedis(esCfg.URL, esCfg.Pass, esCfg.DB, logger)
	defer esClient.Close()

	tracer, tracerCloser := initJaeger(svcName, jCfg.URL, logger)
	defer tracerCloser.Close()

	pubSub, err := mfnats.NewPubSub(natsCfg.URL, svcName, mflogger)
	if err != nil {
		logger.Error("Failed to connect to NATS", zap.Error(err))
		os.Exit(1)
	}
	defer pubSub.Close()

	policiesGRPCConn := connectToGRPC(policiesGRPCCfg, logger)
	defer policiesGRPCConn.Close()

	policiesGRPCTimeout, err := time.ParseDuration(policiesGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", policiesGRPCCfg.Timeout, err.Error())
	}
	policiesGRPCClient := policiesgrpc.NewClient(tracer, policiesGRPCConn, policiesGRPCTimeout)

	fleetGRPCConn := connectToGRPC(fleetGRPCCfg, logger)
	defer fleetGRPCConn.Close()

	fleetGRPCTimeout, err := time.ParseDuration(fleetGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", fleetGRPCCfg.Timeout, err.Error())
	}
	fleetGRPCClient := fleetgrpc.NewClient(tracer, fleetGRPCConn, fleetGRPCTimeout)

	sinksGRPCConn := connectToGRPC(sinksGRPCCfg, logger)
	defer sinksGRPCConn.Close()

	sinksGRPCTimeout, err := time.ParseDuration(sinksGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", sinksGRPCCfg.Timeout, err.Error())
	}
	sinksGRPCClient := sinksgrpc.NewClient(tracer, sinksGRPCConn, sinksGRPCTimeout)

	configRepo, err := config2.NewMemRepo(logger)

	svc := sinker.New(logger, pubSub, esClient, configRepo, policiesGRPCClient, fleetGRPCClient, sinksGRPCClient, promClient)
	defer svc.Stop()

	errs := make(chan error, 2)

	go startHTTPServer(svcCfg.HttpPort, errs, logger)
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
	r.GetFunc("/version", orb.Version(svcName))
	r.Handle("/metrics", promhttp.Handler())
	return r
}

func startHTTPServer(port string, errs chan error, logger *zap.Logger) {
	p := fmt.Sprintf(":%s", port)
	logger.Info("sinker service started, exposed port", zap.String("port", port))
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
