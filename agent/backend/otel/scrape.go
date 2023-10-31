package otel

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confignet"
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
	//var waitGrp sync.WaitGroup
	//waitGrp.Add(1)
	go func() {
		defer execCancelF()
		count := 0
		max := 20
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
						TracerProvider: trace.NewNoopTracerProvider(),
						MeterProvider:  metric.NewMeterProvider(),
						ReportComponentStatus: func(*component.StatusEvent) error {
							return nil
						},
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
				o.logger.Info("Started receiver for OTLP in orb-agent",
					zap.String("host", o.otelReceiverHost), zap.Int("port", o.otelReceiverPort))
				err = receiver.Start(exeCtx, nil)
				if err != nil {
					o.logger.Error("otel receiver startup error", zap.Error(err))
					return
				}
				//waitGrp.Done()
				break
			} else {
				count++
				o.logger.Info("waiting until mqtt client is connected try " + strconv.Itoa(count) + " from " + strconv.Itoa(max))
				time.Sleep(time.Second * time.Duration(count))
				if count >= max {
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
				o.mainCancelFunction()
				return
			}
		}
	}()
	//waitGrp.Wait()
}
