// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package main

import "C"
import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ns1labs/orb/pkg/config"
	natconsume "github.com/ns1labs/orb/pkg/sinks/writer/consumer"
	"github.com/ns1labs/orb/pkg/sinks/writer/prom"
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
	svcName     = "prom-sink"
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

	mflogger, err := mflog.New(os.Stdout, svcCfg.LogLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	esClient := connectToRedis(esCfg.URL, esCfg.Pass, esCfg.DB, mflogger)
	defer esClient.Close()

	pubSub, err := mfnats.NewPubSub(natsCfg.URL, svcName, mflogger)
	if err != nil {
		mflogger.Error(fmt.Sprintf("Failed to connect to NATS: %s", err))
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
			Namespace: "promsink",
			Subsystem: "message_writer",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "promsink",
			Subsystem: "message_writer",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	svc := prom.New(zlog, mflogger, consumerSvc, pubSub, esClient, svcName)

	errs := make(chan error, 2)

	go startHTTPServer(svcCfg.HttpPort, errs, mflogger)
	go svc.Run()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	mflogger.Error(fmt.Sprintf("promsink writer service terminated: %s", err))
}

func startHTTPServer(port string, errs chan error, logger mflog.Logger) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("promsink writer service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, mfwriters.MakeHandler(svcName))
}

func connectToRedis(URL, pass string, cacheDB string, logger mflog.Logger) *redis.Client {
	db, err := strconv.Atoi(cacheDB)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to cache: %s", err))
		os.Exit(1)
	}

	return redis.NewClient(&redis.Options{
		Addr:     URL,
		Password: pass,
		DB:       db,
	})
}
