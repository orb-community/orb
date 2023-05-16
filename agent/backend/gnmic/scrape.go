/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package gnmic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"github.com/orb-community/orb/agent/otel"
	"github.com/orb-community/orb/agent/otel/otlpmqttexporter"
	configutil "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	typeStr            = "orb-gnmic"
	PromScheme         = "http"
	PromScrapeInterval = 60 * time.Second
	PromScrapeTimeout  = 15 * time.Second
	PromMetricsPath    = "/metrics"
)

func (d *gnmicBackend) createOtlpMqttLogsExporter(ctx context.Context, cancelFunc context.CancelFunc) (exporter.Logs, error) {

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

func (d *gnmicBackend) createOtlpMqttMetricsExporter(ctx context.Context, cancelFunc context.CancelFunc) (exporter.Metrics, error) {

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

func (d *gnmicBackend) scrapeOtlp(ctx context.Context) {
	exeCtx, execCancelF := context.WithCancel(d.ctx)
	policyID := ctx.Value("policy_id").(string)
	policyName := ctx.Value("policy_name").(string)
	// add policyName in context
	go func() {
		defer execCancelF()
		var err error
		count := 0
		for {
			if d.mqttClient != nil {
				// Metrics
				d.exporter[policyID], err = d.createOtlpMqttMetricsExporter(exeCtx, execCancelF)
				if err != nil {
					d.logger.Error("failed to create a metrics exporter", zap.Error(err))
					return
				}
				promFactory := prometheusreceiver.NewFactory()
				promCfg := promFactory.CreateDefaultConfig().(*prometheusreceiver.Config)
				httpConfig := configutil.HTTPClientConfig{}
				scrapeConfig := &config.ScrapeConfig{
					ScrapeInterval:  model.Duration(PromScrapeInterval),
					ScrapeTimeout:   model.Duration(PromScrapeTimeout),
					JobName:         fmt.Sprintf("%s-%s", typeStr, policyName),
					HonorTimestamps: true,
					Scheme:          PromScheme,

					MetricsPath: PromMetricsPath,
					Params:      d.PromParams,
					ServiceDiscoveryConfigs: discovery.Configs{
						&discovery.StaticConfig{
							{
								Targets: []model.LabelSet{
									{model.AddressLabel: model.LabelValue(d.otelMetricsReceiverHost)},
								},
							},
						},
					},
				}
				scrapeConfig.HTTPClientConfig = httpConfig
				promCfg.PrometheusConfig = &config.Config{ScrapeConfigs: []*config.ScrapeConfig{
					scrapeConfig,
				}}

				set := receiver.CreateSettings{
					TelemetrySettings: component.TelemetrySettings{
						Logger:         d.logger,
						TracerProvider: trace.NewNoopTracerProvider(),
						MeterProvider:  global.MeterProvider(),
					},
					BuildInfo: component.NewDefaultBuildInfo(),
				}
				d.receiver[policyID], err = promFactory.CreateMetricsReceiver(exeCtx, set, promCfg, d.exporter[policyID])
				if err != nil {
					d.logger.Error("failed to create a metrics receiver", zap.Error(err))
					return
				}

				err = d.exporter[policyID].Start(exeCtx, nil)
				if err != nil {
					d.logger.Error("otel mqtt metrics exporter startup error", zap.Error(err))
					return
				}

				err = d.receiver[policyID].Start(exeCtx, nil)
				if err != nil {
					d.logger.Error("otel receiver metrics startup error", zap.Error(err))
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
				err := d.exporter[policyID].Shutdown(exeCtx)
				if err != nil {
					return
				}
				err = d.receiver[policyID].Shutdown(exeCtx)
				if err != nil {
					return
				}
				d.logger.Info("stopped gnmic telemetry collector policy:" + policyID)
				return
			}
		}
	}()
}
