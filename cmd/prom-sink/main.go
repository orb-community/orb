// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package main

import "C"
import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mainflux/mainflux/bootstrap/redis/consumer"
	sinkwriter "github.com/ns1labs/orb/pkg/sinks/writer"
	"github.com/ns1labs/orb/pkg/sinks/writer/prom"
	natconsume "github.com/ns1labs/orb/pkg/sinks/writer/prom/consumer"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/mainflux/mainflux/consumers/writers/api"
	mflog "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging/nats"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	svcName = "prom-sink"
)

type config struct {
	NatsURL string `mapstructure:"nats_url"`

	EsURL  string `mapstructure:"es_url"`
	EsPass string `mapstructure:"es_pass"`
	EsDB   string `mapstructure:"es_db"`

	LogLevel   string `mapstructure:"log_level"`
	Port       string `mapstructure:"port"`
	ConfigPath string `mapstructure:"config_path"`

	// todo sinks gRPC
	// todo fleet mgr gRPC
}

func loadConfig() config {

	mfC := viper.New()
	mfC.SetEnvPrefix("mf")

	mfC.SetDefault("nats_url", "nats://localhost:4222")

	mfC.SetDefault("es_url", "localhost:6379")
	mfC.SetDefault("es_pass", "")
	mfC.SetDefault("es_db", "0")

	mfC.AutomaticEnv()

	orbC := viper.New()
	orbC.SetEnvPrefix("orb_prom_sink")

	orbC.SetDefault("config_path", "/config.toml")
	orbC.SetDefault("log_level", "error")
	orbC.SetDefault("port", "8190")

	orbC.AutomaticEnv()

	var c config
	mfC.Unmarshal(&c)
	orbC.Unmarshal(&c)
	return c

}

func main() {

	cfg := loadConfig()

	logger, err := mflog.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	esClient := connectToRedis(cfg.EsURL, cfg.EsPass, cfg.EsDB, logger)
	defer esClient.Close()

	pubSub, err := nats.NewPubSub(cfg.NatsURL, svcName, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to NATS: %s", err))
		os.Exit(1)
	}
	defer pubSub.Close()

	// todo fleet grpc
	// todo sink grpc

	svc := newService(logger)

	errs := make(chan error, 2)

	go startHTTPServer(cfg.Port, errs, logger)
	go subscribeToOrbES(svc)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("promsink writer service terminated: %s", err))
}

func newService(logger mflog.Logger) sinkwriter.Service {
	zlog, _ := zap.NewProduction()
	consumerSvc := natconsume.New(zlog)
	consumerSvc = api.LoggingMiddleware(consumerSvc, logger)
	consumerSvc = api.MetricsMiddleware(
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
	svc := prom.New()
	return svc
}

func startHTTPServer(port string, errs chan error, logger mflog.Logger) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("promsink writer service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHandler(svcName))
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

func subscribeToOrbES(svc sinkwriter.Service, client *redis.Client, esconsumer string, logger mflog.Logger) {
	eventStore := consumer.NewEventStore(svc, client, esconsumer, logger)
	logger.Info("Subscribed to Redis Event Store")
	if err := eventStore.Subscribe("orb.policy"); err != nil {
		logger.Warn(fmt.Sprintf("orb prometheus sync service failed to subscribe to event sourcing: %s", err))
	}
}
