// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package main

import "C"
import (
	"fmt"
	"github.com/ns1labs/orb/pkg/mainflux/consumers/writers/promsink"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/mainflux/mainflux/consumers"
	"github.com/mainflux/mainflux/consumers/writers/api"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/ns1labs/orb/pkg/mainflux/transformers/passthrough"
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
	orbC.SetDefault("port", "8180")

	orbC.AutomaticEnv()

	var c config
	mfC.Unmarshal(&c)
	orbC.Unmarshal(&c)
	return c

}

func main() {

	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	pubSub, err := nats.NewPubSub(cfg.NatsURL, "", logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to NATS: %s", err))
		os.Exit(1)
	}
	defer pubSub.Close()

	// prometheus connection: https://github.com/timescale/promscale/blob/master/docs/writing_to_promscale.md
	//db := connectToDB(cfg.dbConfig, logger)
	//defer db.Close()

	repo := newService( /*db, */ logger)
	t := passthrough.New()

	if err = consumers.Start(pubSub, repo, t, cfg.ConfigPath, logger); err != nil {
		logger.Error(fmt.Sprintf("Failed to create promsink writer: %s", err))
	}

	errs := make(chan error, 2)

	go startHTTPServer(cfg.Port, errs, logger)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("promsink writer service terminated: %s", err))
}

//func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sqlx.DB {
//	db, err := postgres.Connect(dbConfig)
//	if err != nil {
//		logger.Error(fmt.Sprintf("Failed to connect to Postgres: %s", err))
//		os.Exit(1)
//	}
//	return db
//}

func newService( /*db *sqlx.DB, */ logger logger.Logger) consumers.Consumer {
	zlog, _ := zap.NewProduction()
	svc := promsink.New(zlog)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
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

	return svc
}

func startHTTPServer(port string, errs chan error, logger logger.Logger) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("promsink writer service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHandler(svcName))
}
