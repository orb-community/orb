package otel

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func (o *openTelemetryBackend) receiveOtlp() {
	exeCtx, execCancelF := context.WithCancel(o.mainContext)
	go func() {
		defer execCancelF()
		count := 0
		for {
			if o.mqttClient != nil {
				exporter, err := o.createOtlpMqttExporter(exeCtx, execCancelF)
				if err != nil {
					o.logger.Error("failed to create a exporter", zap.Error(err))
					return
				}
				pFactory := otlpreceiver.NewFactory()
				cfg := pFactory.CreateDefaultConfig()
				cfg.(*otlpreceiver.Config).Protocols = otlpreceiver.Protocols{
					HTTP: &otlpreceiver.HTTPConfig{
						HTTPServerSettings: &confighttp.HTTPServerSettings{
							Endpoint: o.otelReceiverHost + ":" + strconv.Itoa(o.otelReceiverPort),
						},
					},
				}
				set := receiver.CreateSettings{
					TelemetrySettings: component.TelemetrySettings{
						Logger:         o.logger,
						TracerProvider: trace.NewNoopTracerProvider(),
						MeterProvider:  metric.NewMeterProvider(),
					},
					BuildInfo: component.NewDefaultBuildInfo(),
				}
				receiver, err := pFactory.CreateMetricsReceiver(exeCtx, set, cfg, exporter)
				if err != nil {
					o.logger.Error("failed to create a receiver", zap.Error(err))
					return
				}
				err = exporter.Start(exeCtx, nil)
				if err != nil {
					o.logger.Error("otel mqtt exporter startup error", zap.Error(err))
					return
				}
				err = receiver.Start(exeCtx, nil)
				if err != nil {
					o.logger.Error("otel receiver startup error", zap.Error(err))
					return
				}
				break
			} else {
				count++
				o.logger.Info("waiting until mqtt client is connected try " + strconv.Itoa(count) + " from 10")
				time.Sleep(time.Second * 3)
				if count >= 10 {
					execCancelF()
					o.mainCancelFunction()
					break
				}
			}
		}
		for {
			select {
			case <-exeCtx.Done():
				o.mainContext.Done()
				o.mainCancelFunction()
			case <-o.mainContext.Done():
				o.logger.Info("stopped Orb OpenTelemetry agent collector")
				return
			}
		}
	}()
}
