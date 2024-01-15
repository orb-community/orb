package otel

import (
	"context"
	"errors"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func (o *openTelemetryBackend) receiveOtlp() {
	exeCtx, execCancelF := context.WithCancelCause(o.mainContext)
	go func() {
		count := 0
		maxRetries := 20
		for {
			if o.mqttClient != nil {
				if ok := o.startOtelMetric(exeCtx, execCancelF); !ok {
					o.logger.Error("failed to start otel metric")
					return
				}
				//if o.startOtelTraces(exeCtx, execCancelF) {
				//	return
				//}
				//if ok := o.startOtelLogs(exeCtx, execCancelF); !ok {
				//	return
				//}
				o.logger.Info("started otel receiver for opentelemetry")
				break
			} else {
				count++
				o.logger.Info("waiting until mqtt client is connected try " + strconv.Itoa(count) + " from " + strconv.Itoa(maxRetries))
				time.Sleep(time.Second * time.Duration(count))
				if count >= maxRetries {
					execCancelF(errors.New("mqtt client is not connected"))
					o.mainCancelFunction()
					break
				}
			}
		}
		for {
			select {
			case <-exeCtx.Done():
				o.logger.Info("stopped receiver context, pktvisor will not scrape metrics", zap.Error(context.Cause(exeCtx)))
				o.mainContext.Done()
				o.mainCancelFunction()
			case <-o.mainContext.Done():
				o.logger.Info("stopped Orb OpenTelemetry agent collector")
				o.mainCancelFunction()
				return
			}
		}
	}()
}

func (o *openTelemetryBackend) startOtelMetric(exeCtx context.Context, execCancelF context.CancelCauseFunc) bool {
	var err error
	o.metricsExporter, err = o.createOtlpMetricMqttExporter(exeCtx, execCancelF)
	if err != nil {
		o.logger.Error("failed to create a exporter", zap.Error(err))
		return false
	}
	pFactory := otlpreceiver.NewFactory()
	cfg := pFactory.CreateDefaultConfig()
	cfg.(*otlpreceiver.Config).Protocols = otlpreceiver.Protocols{
		GRPC: &configgrpc.GRPCServerSettings{
			NetAddr: confignet.NetAddr{
				Endpoint:  o.otelReceiverHost + ":" + strconv.Itoa(o.otelReceiverPort),
				Transport: "tcp",
			},
		},
	}
	set := receiver.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         o.logger,
			TracerProvider: noop.NewTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			ReportComponentStatus: func(*component.StatusEvent) error {
				return nil
			},
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
	o.metricsReceiver, err = pFactory.CreateMetricsReceiver(exeCtx, set, cfg, o.metricsExporter)
	if err != nil {
		o.logger.Error("failed to create a receiver", zap.Error(err))
		return false
	}
	err = o.metricsExporter.Start(exeCtx, nil)
	if err != nil {
		o.logger.Error("otel mqtt exporter startup error", zap.Error(err))
		return false
	}
	o.logger.Info("Started receiver for OTLP in orb-agent",
		zap.String("host", o.otelReceiverHost), zap.Int("port", o.otelReceiverPort))
	err = o.metricsReceiver.Start(exeCtx, nil)
	if err != nil {
		o.logger.Error("otel receiver startup error", zap.Error(err))
		return false
	}
	return true
}

// TODO fix when create otlpmqtt trace
//func (o *openTelemetryBackend) startOtelTraces(exeCtx context.Context, execCancelF context.CancelFunc) bool {
//	if o.tracesExporter != nil {
//		return true
//	}
//	var err error
//	o.tracesExporter, err = o.createOtlpTraceMqttExporter(exeCtx, execCancelF)
//	if err != nil {
//		o.logger.Error("failed to create a exporter", zap.Error(err))
//		return true
//	}
//	pFactory := otlpreceiver.NewFactory()
//	cfg := pFactory.CreateDefaultConfig()
//	cfg.(*otlpreceiver.Config).Protocols = otlpreceiver.Protocols{
//		GRPC: &configgrpc.GRPCServerSettings{
//			NetAddr: confignet.NetAddr{
//				Endpoint:  o.otelReceiverHost + ":" + strconv.Itoa(o.otelReceiverPort),
//				Transport: "tcp",
//			},
//		},
//	}
//	set := receiver.CreateSettings{
//		TelemetrySettings: component.TelemetrySettings{
//			Logger:         o.logger,
//			TracerProvider: trace.NewNoopTracerProvider(),
//			MeterProvider:  metric.NewMeterProvider(),
//			ReportComponentStatus: func(*component.StatusEvent) error {
//				return nil
//			},
//		},
//		BuildInfo: component.NewDefaultBuildInfo(),
//	}
//	o.tracesReceiver, err = pFactory.CreateTracesReceiver(exeCtx, set, cfg, o.tracesExporter)
//	if err != nil {
//		o.logger.Error("failed to create a receiver", zap.Error(err))
//		return true
//	}
//	err = o.metricsExporter.Start(exeCtx, nil)
//	if err != nil {
//		o.logger.Error("otel mqtt exporter startup error", zap.Error(err))
//		return true
//	}
//	o.logger.Info("Started receiver for OTLP in orb-agent",
//		zap.String("host", o.otelReceiverHost), zap.Int("port", o.otelReceiverPort))
//	err = o.metricsReceiver.Start(exeCtx, nil)
//	if err != nil {
//		o.logger.Error("otel receiver startup error", zap.Error(err))
//		return true
//	}
//	return false
//}
//
//func (o *openTelemetryBackend) startOtelLogs(exeCtx context.Context, execCancelF context.CancelFunc) bool {
//	if o.logsExporter != nil {
//		return true
//	}
//	var err error
//	o.logsExporter, err = o.createOtlpLogsMqttExporter(exeCtx, execCancelF)
//	if err != nil {
//		o.logger.Error("failed to create a exporter", zap.Error(err))
//		return false
//	}
//	pFactory := otlpreceiver.NewFactory()
//	cfg := pFactory.CreateDefaultConfig()
//	cfg.(*otlpreceiver.Config).Protocols = otlpreceiver.Protocols{
//		GRPC: &configgrpc.GRPCServerSettings{
//			NetAddr: confignet.NetAddr{
//				Endpoint:  o.otelReceiverHost + ":" + strconv.Itoa(o.otelReceiverPort),
//				Transport: "tcp",
//			},
//		},
//	}
//	set := receiver.CreateSettings{
//		TelemetrySettings: component.TelemetrySettings{
//			Logger:         o.logger,
//			TracerProvider: trace.NewNoopTracerProvider(),
//			MeterProvider:  metric.NewMeterProvider(),
//			ReportComponentStatus: func(*component.StatusEvent) error {
//				return nil
//			},
//		},
//		BuildInfo: component.NewDefaultBuildInfo(),
//	}
//	o.metricsReceiver, err = pFactory.CreateLogsReceiver(exeCtx, set, cfg, o.logsExporter)
//	if err != nil {
//		o.logger.Error("failed to create a receiver", zap.Error(err))
//		return false
//	}
//	err = o.metricsExporter.Start(exeCtx, nil)
//	if err != nil {
//		o.logger.Error("otel mqtt exporter startup error", zap.Error(err))
//		return false
//	}
//	o.logger.Info("Started receiver for OTLP in orb-agent",
//		zap.String("host", o.otelReceiverHost), zap.Int("port", o.otelReceiverPort))
//	err = o.metricsReceiver.Start(exeCtx, nil)
//	if err != nil {
//		o.logger.Error("otel receiver startup error", zap.Error(err))
//		return false
//	}
//	return true
//}
