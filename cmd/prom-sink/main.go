// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-zoo/bone"
	"github.com/ns1labs/orb"
	"github.com/ns1labs/orb/pkg/config"
	natconsume "github.com/ns1labs/orb/sinks/writer/consumer"
	"github.com/ns1labs/orb/sinks/writer/prom"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	mfwriters "github.com/mainflux/mainflux/consumers/writers/api"
	mflog "github.com/mainflux/mainflux/logger"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	svcName     = "prom_sink"
	fullSvcName = "orb-prom-sink"
	envPrefix   = "orb_prom_sink"
	httpPort    = "8201"
)

func main() {

	natsCfg := config.LoadNatsConfig(envPrefix)
	esCfg := config.LoadEsConfig(envPrefix)
	svcCfg := config.LoadBaseServiceConfig(envPrefix, httpPort)

	// todo sinks gRPC
	// todo fleet mgr gRPC

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

	pubSub, err := mfnats.NewPubSub(natsCfg.URL, svcName, mflogger)
	if err != nil {
		logger.Error("Failed to connect to NATS", zap.Error(err))
		os.Exit(1)
	}
	defer pubSub.Close()

	// todo fleet grpc
	// todo sink grpc

	zlog, _ := zap.NewProduction()
	consumerSvc := natconsume.New(zlog)
	consumerSvc = mfwriters.LoggingMiddleware(consumerSvc, mflogger)
	consumerSvc = mfwriters.MetricsMiddleware(
		consumerSvc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: svcName,
			Subsystem: "message_writer",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: svcName,
			Subsystem: "message_writer",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	svc := prom.New(zlog, mflogger, consumerSvc, pubSub, esClient, svcName)

	errs := make(chan error, 2)

	go startHTTPServer(svcCfg.HttpPort, errs, logger)
	go svc.Run()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error("promsink writer service terminated", zap.Error(err))
}

func makeHandler(svcName string) http.Handler {
	r := bone.New()
	r.GetFunc("/version", orb.Version(svcName))
	r.Handle("/metrics", promhttp.Handler())
	return r
}

func startHTTPServer(port string, errs chan error, logger *zap.Logger) {
	p := fmt.Sprintf(":%s", port)
	logger.Info("promsink writer service started, exposed port", zap.String("port", port))
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
