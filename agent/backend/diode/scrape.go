/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package diode

import (
	"context"
	"go.opentelemetry.io/otel/sdk/metric"
	"strconv"
	"time"

	"github.com/orb-community/orb/agent/otel"
	"github.com/orb-community/orb/agent/otel/otlpmqttexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	otlpProtocol = "tcp"
)

func (d *diodeBackend) createOtlpMqttExporter(ctx context.Context, cancelFunc context.CancelFunc) (exporter.Logs, error) {

	bridgeService := otel.NewBridgeService(ctx, cancelFunc, &d.policyRepo, d.agentTags)
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

func (d *diodeBackend) receiveOtlp() {
	exeCtx, execCancelF := context.WithCancel(d.ctx)
	go func() {
		defer execCancelF()
		var err error
		count := 0
		for {
			if d.mqttClient != nil {
				d.exporter, err = d.createOtlpMqttExporter(exeCtx, execCancelF)
				if err != nil {
					d.logger.Error("failed to create a exporter", zap.Error(err))
					return
				}
				pFactory := otlpreceiver.NewFactory()
				cfg := pFactory.CreateDefaultConfig().(*otlpreceiver.Config)
				cfg.HTTP = nil
				cfg.GRPC.NetAddr.Endpoint = d.otelReceiverHost + ":" + strconv.Itoa(d.otelReceiverPort)
				cfg.GRPC.NetAddr.Transport = otlpProtocol

				set := receiver.CreateSettings{
					TelemetrySettings: component.TelemetrySettings{
						Logger:         d.logger,
						TracerProvider: trace.NewNoopTracerProvider(),
						MeterProvider:  metric.NewMeterProvider(),
					},
					BuildInfo: component.NewDefaultBuildInfo(),
				}
				d.receiver, err = pFactory.CreateLogsReceiver(exeCtx, set, cfg, d.exporter)
				if err != nil {
					d.logger.Error("failed to create a receiver", zap.Error(err))
					return
				}

				err = d.exporter.Start(exeCtx, nil)
				if err != nil {
					d.logger.Error("otel mqtt exporter startup error", zap.Error(err))
					return
				}

				err = d.receiver.Start(exeCtx, nil)
				if err != nil {
					d.logger.Error("otel receiver startup error", zap.Error(err))
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
				err := d.exporter.Shutdown(exeCtx)
				if err != nil {
					return
				}
				err = d.receiver.Shutdown(exeCtx)
				if err != nil {
					return
				}
				d.logger.Info("stopped diode agent OpenTelemetry collector")
				return
			}
		}
	}()
}
