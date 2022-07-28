// Copyright 2019 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Program otelcontribcol is an extension to the OpenTelemetry Collector
// that includes additional components, some vendor-specific, contributed
// from the wider community.

//go:build !testbed
// +build !testbed

package main

import (
	"context"
	"github.com/ns1labs/orb/otelcollector"
	"github.com/ns1labs/orb/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"runtime"
	"strings"
)

const (
	svcName   = "otelcollector"
	envPrefix = "orb_otelcollector"
	httpPort  = "8204"
)

var logger *zap.Logger

func main() {

	atomicLevel := zap.NewAtomicLevel()
	svcCfg := config.LoadBaseServiceConfig(envPrefix, httpPort)
	sinkerGrpcCfg := config.LoadGRPCConfig(envPrefix, "sinker")
	policiesGrpcCfg := config.LoadGRPCConfig(envPrefix, "policies")
	sinksGrpcCfg := config.LoadGRPCConfig(envPrefix, "sinks")

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
	logger.Info("initializing logger")
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	mainCtx, cancel := controlContext(logger)
	defer cancel()

	otelcollector.StartCollector(mainCtx, *logger, svcCfg, sinkerGrpcCfg, policiesGrpcCfg, sinksGrpcCfg)
}

// Controls context cancellation
func controlContext(logger *zap.Logger) (context.Context, func()) {
	ctx, cancel := context.WithCancel(config.ContextWithLogger(context.Background(), logger))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
			return
		}
		logger.Debug("exiting with goroutines and gocalls", zap.Int("goroutines", runtime.NumGoroutine()), zap.Int64("gocalls", runtime.NumCgoCall()))
	}()

	return ctx, func() {
		signal.Stop(c)
		cancel()
	}
}
