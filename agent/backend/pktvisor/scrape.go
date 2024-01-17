/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/trace/noop"
	"net/http"
	"strconv"
	"time"

	"github.com/orb-community/orb/agent/otel"
	"github.com/orb-community/orb/agent/otel/otlpmqttexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
)

func (p *pktvisorBackend) scrapeMetrics(period uint) (map[string]interface{}, error) {
	var metrics map[string]interface{}
	err := p.request(fmt.Sprintf("policies/__all/metrics/bucket/%d", period), &metrics, http.MethodGet, http.NoBody, "application/json", ScrapeTimeout)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (p *pktvisorBackend) createOtlpMqttExporter(ctx context.Context, cancelFunc context.CancelCauseFunc) (exporter.Metrics, error) {
	bridgeService := otel.NewBridgeService(ctx, cancelFunc, &p.policyRepo, p.agentTags)
	var cfg component.Config
	if p.mqttClient != nil {
		cfg = otlpmqttexporter.CreateConfigClient(p.mqttClient, p.otlpMetricsTopic, p.pktvisorVersion, bridgeService)
	} else {
		cfg = otlpmqttexporter.CreateConfig(p.mqttConfig.Address, p.mqttConfig.Id, p.mqttConfig.Key,
			p.mqttConfig.ChannelID, p.pktvisorVersion, p.otlpMetricsTopic, bridgeService)
	}

	set := otlpmqttexporter.CreateDefaultSettings(p.logger)
	// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
	return otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
}

func (p *pktvisorBackend) startOtelMetric(exeCtx context.Context, execCancelF context.CancelCauseFunc) bool {
	var err error
	p.exporter, err = p.createOtlpMqttExporter(exeCtx, execCancelF)
	if err != nil {
		p.logger.Error("failed to create a exporter", zap.Error(err))
		return false
	}
	pFactory := otlpreceiver.NewFactory()
	cfg := pFactory.CreateDefaultConfig()
	cfg.(*otlpreceiver.Config).Protocols = otlpreceiver.Protocols{
		HTTP: &otlpreceiver.HTTPConfig{
			HTTPServerSettings: &confighttp.HTTPServerSettings{
				Endpoint: p.otelReceiverHost + ":" + strconv.Itoa(p.otelReceiverPort),
			},
			MetricsURLPath: "/v1/metrics",
		},
	}
	set := receiver.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         p.logger,
			TracerProvider: noop.NewTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
			ReportComponentStatus: func(*component.StatusEvent) error {
				return nil
			},
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}

	p.receiver, err = pFactory.CreateMetricsReceiver(exeCtx, set, cfg, p.exporter)
	if err != nil {
		p.logger.Error("failed to create a receiver", zap.Error(err))
		return false
	}
	err = p.exporter.Start(exeCtx, nil)
	if err != nil {
		p.logger.Error("otel mqtt exporter startup error", zap.Error(err))
		return false
	}
	p.logger.Info("Started receiver for OTLP in orb-agent",
		zap.String("host", p.otelReceiverHost), zap.Int("port", p.otelReceiverPort))
	err = p.receiver.Start(exeCtx, nil)
	if err != nil {
		p.logger.Error("otel receiver startup error", zap.Error(err))
		return false
	}
	return true
}

func (p *pktvisorBackend) receiveOtlp() {
	exeCtx, execCancelF := context.WithCancelCause(p.ctx)
	go func() {
		count := 0
		for {
			if p.mqttClient != nil {
				if ok := p.startOtelMetric(exeCtx, execCancelF); !ok {
					p.logger.Error("failed to start otel metric")
					return
				}
				p.logger.Info("started otel receiver for pktvisor")
				break
			} else {
				count++
				p.logger.Info("waiting until mqtt client is connected try " + strconv.Itoa(count) + " from 10")
				time.Sleep(time.Second * time.Duration(count))
				if count >= 10 {
					execCancelF(errors.New("mqtt client is not connected"))
					_ = p.Stop(exeCtx)
					break
				}
			}
		}
		for {
			select {
			case <-exeCtx.Done():
				p.logger.Info("stopped receiver context, pktvisor will not scrape metrics", zap.Error(context.Cause(exeCtx)))
				p.cancelFunc()
				_ = p.exporter.Shutdown(exeCtx)
				_ = p.receiver.Shutdown(exeCtx)
			case <-p.ctx.Done():
				p.logger.Info("stopped pktvisor main context, stopping receiver")
				execCancelF(errors.New("stopped pktvisor main context"))
				_ = p.exporter.Shutdown(exeCtx)
				_ = p.receiver.Shutdown(exeCtx)
				return
			}
		}
	}()
}
