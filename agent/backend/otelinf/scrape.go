/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otelinf

import (
	"context"
	"strconv"
	"time"

	"github.com/orb-community/orb/agent/otel"
	"github.com/orb-community/orb/agent/otel/otlpmqttexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	otlpProtocol = "tcp"
)

func (d *otelinfBackend) createOtlpMqttLogsExporter(ctx context.Context, cancelFunc context.CancelFunc) (exporter.Logs, error) {

	bridgeService := otel.NewBridgeService(ctx, &d.policyRepo, d.agentTags)
	if d.mqttClient != nil {
		cfg := otlpmqttexporter.CreateConfigClient(d.mqttClient, d.logTopic, d.version, bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(d.logger)
		// Create the OTLP logs exporter that'll receive and verify the logs produced.
		exporter, err := otlpmqttexporter.CreateLogsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	} else {
		cfg := otlpmqttexporter.CreateConfig(d.mqttConfig.Address, d.mqttConfig.Id, d.mqttConfig.Key,
			d.mqttConfig.ChannelID, d.version, d.logTopic, bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(d.logger)
		// Create the OTLP logs exporter that'll receive and verify the logs produced.
		exporter, err := otlpmqttexporter.CreateLogsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	}

}

func (d *otelinfBackend) createOtlpMqttMetricsExporter(ctx context.Context, cancelFunc context.CancelFunc) (exporter.Metrics, error) {

	bridgeService := otel.NewBridgeService(ctx, &d.policyRepo, d.agentTags)
	if d.mqttClient != nil {
		cfg := otlpmqttexporter.CreateConfigClient(d.mqttClient, d.metricTopic, d.version, bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(d.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		exporter, err := otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	} else {
		cfg := otlpmqttexporter.CreateConfig(d.mqttConfig.Address, d.mqttConfig.Id, d.mqttConfig.Key,
			d.mqttConfig.ChannelID, d.version, d.metricTopic, bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(d.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		exporter, err := otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	}

}

func (d *otelinfBackend) receiveOtlp() {
	exeCtx, execCancelF := context.WithCancel(d.ctx)
	go func() {
		defer execCancelF()
		var err error
		count := 0
		for {
			if d.mqttClient != nil {
				// Metrics
				d.exporterMetrics, err = d.createOtlpMqttMetricsExporter(exeCtx, execCancelF)
				if err != nil {
					d.logger.Error("failed to create a metrics exporter", zap.Error(err))
					return
				}
				pFactory := otlpreceiver.NewFactory()
				cfg := pFactory.CreateDefaultConfig().(*otlpreceiver.Config)
				cfg.HTTP = nil
				cfg.GRPC.NetAddr.Endpoint = d.otelMetricsReceiverHost + ":" + strconv.Itoa(d.otelMetricsReceiverPort)
				cfg.GRPC.NetAddr.Transport = otlpProtocol

				set := receiver.CreateSettings{
					TelemetrySettings: component.TelemetrySettings{
						Logger:         d.logger,
						TracerProvider: trace.NewNoopTracerProvider(),
						MeterProvider:  global.MeterProvider(),
					},
					BuildInfo: component.NewDefaultBuildInfo(),
				}
				d.receiverMetrics, err = pFactory.CreateMetricsReceiver(exeCtx, set, cfg, d.exporterMetrics)
				if err != nil {
					d.logger.Error("failed to create a metrics receiver", zap.Error(err))
					return
				}

				err = d.exporterMetrics.Start(exeCtx, nil)
				if err != nil {
					d.logger.Error("otel mqtt metrics exporter startup error", zap.Error(err))
					return
				}

				err = d.receiverMetrics.Start(exeCtx, nil)
				if err != nil {
					d.logger.Error("otel receiver metrics startup error", zap.Error(err))
					return
				}
				// Logs
				d.exporterLogs, err = d.createOtlpMqttLogsExporter(exeCtx, execCancelF)
				if err != nil {
					d.logger.Error("failed to create a logs exporter", zap.Error(err))
					return
				}
				pFactoryLogs := otlpreceiver.NewFactory()
				cfgLogs := pFactoryLogs.CreateDefaultConfig().(*otlpreceiver.Config)
				cfgLogs.HTTP = nil
				cfgLogs.GRPC.NetAddr.Endpoint = d.otelLogsReceiverHost + ":" + strconv.Itoa(d.otelLogsReceiverPort)
				cfgLogs.GRPC.NetAddr.Transport = otlpProtocol

				setLogs := receiver.CreateSettings{
					TelemetrySettings: component.TelemetrySettings{
						Logger:         d.logger,
						TracerProvider: trace.NewNoopTracerProvider(),
						MeterProvider:  global.MeterProvider(),
					},
					BuildInfo: component.NewDefaultBuildInfo(),
				}
				d.receiverLogs, err = pFactory.CreateLogsReceiver(exeCtx, setLogs, cfgLogs, d.exporterLogs)
				if err != nil {
					d.logger.Error("failed to create a logs receiver", zap.Error(err))
					return
				}

				err = d.exporterLogs.Start(exeCtx, nil)
				if err != nil {
					d.logger.Error("otel mqtt logs exporter startup error", zap.Error(err))
					return
				}

				err = d.receiverLogs.Start(exeCtx, nil)
				if err != nil {
					d.logger.Error("otel logs receiver startup error", zap.Error(err))
					return
				}
				break
			} else {
				count++
				d.logger.Info("waiting until mqtt client is connected try " + strconv.Itoa(count) + " from 10")
				time.Sleep(time.Second * 3)
				if count >= 10 {
					execCancelF()
					_ = d.Stop(exeCtx)
					break
				}
			}
		}
		for {
			select {
			case <-exeCtx.Done():
				d.ctx.Done()
				d.cancelFunc()
			case <-d.ctx.Done():
				err := d.exporterMetrics.Shutdown(exeCtx)
				if err != nil {
					return
				}
				err = d.receiverMetrics.Shutdown(exeCtx)
				if err != nil {
					return
				}
				err = d.exporterLogs.Shutdown(exeCtx)
				if err != nil {
					return
				}
				err = d.receiverLogs.Shutdown(exeCtx)
				if err != nil {
					return
				}
				d.logger.Info("stopped otelinf agent OpenTelemetry collector")
				return
			}
		}
	}()
}
