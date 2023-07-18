/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/orb-community/orb/agent/otel"

	"github.com/orb-community/orb/agent/otel/otlpmqttexporter"
	"github.com/orb-community/orb/fleet"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
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

func (p *pktvisorBackend) createOtlpMqttExporter(ctx context.Context, cancelFunc context.CancelFunc) (exporter.Metrics, error) {

	bridgeService := otel.NewBridgeService(ctx, cancelFunc, &p.policyRepo, p.agentTags)
	if p.mqttClient != nil {
		cfg := otlpmqttexporter.CreateConfigClient(p.mqttClient, p.otlpMetricsTopic, p.pktvisorVersion, bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(p.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		metricsExporter, err := otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return metricsExporter, nil
	} else {
		cfg := otlpmqttexporter.CreateConfig(p.mqttConfig.Address, p.mqttConfig.Id, p.mqttConfig.Key,
			p.mqttConfig.ChannelID, p.pktvisorVersion, p.otlpMetricsTopic, bridgeService)
		set := otlpmqttexporter.CreateDefaultSettings(p.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		metricsExporter, err := otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return metricsExporter, nil
	}

}

func (p *pktvisorBackend) scrapeDefault() error {
	// scrape all policy json output with one call every minute.
	// TODO support policies with custom bucket times
	job, err := p.scraper.Every(1).Minute().WaitForSchedule().Do(func() {
		metrics, err := p.scrapeMetrics(1)
		if err != nil {
			p.logger.Error("scrape failed", zap.Error(err))
			return
		}
		if len(metrics) == 0 {
			p.logger.Warn("scrape: no policies found, skipping")
			return
		}

		var batchPayload []fleet.AgentMetricsRPCPayload
		totalSize := 0
		for pName, pMetrics := range metrics {
			policyData, err := p.policyRepo.GetByName(pName)
			if err != nil {
				p.logger.Warn("skipping pktvisor policy not managed by orb", zap.String("policy", pName), zap.Error(err))
				continue
			}
			payloadData, err := json.Marshal(pMetrics)
			if err != nil {
				p.logger.Error("error marshalling scraped metric json", zap.String("policy", pName), zap.Error(err))
				continue
			}
			metricPayload := fleet.AgentMetricsRPCPayload{
				PolicyID:   policyData.ID,
				PolicyName: policyData.Name,
				Datasets:   policyData.GetDatasetIDs(),
				Format:     "json",
				BEVersion:  p.pktvisorVersion,
				Data:       payloadData,
			}
			batchPayload = append(batchPayload, metricPayload)
			totalSize += len(payloadData)
			policyData.LastScrapeBytes = int64(totalSize)
			policyData.LastScrapeTS = time.Now()
			err = p.policyRepo.Update(policyData)
			if err != nil {
				p.logger.Error("unable to update policy repo during scrape", zap.Error(err))
			}
			p.logger.Info("scraped metrics for policy", zap.String("policy", pName), zap.String("policy_id", policyData.ID), zap.Int("payload_size_b", len(payloadData)))
		}

		rpc := fleet.AgentMetricsRPC{
			SchemaVersion: fleet.CurrentRPCSchemaVersion,
			Func:          fleet.AgentMetricsRPCFunc,
			Payload:       batchPayload,
		}

		body, err := json.Marshal(rpc)
		if err != nil {
			p.logger.Error("error marshalling metric rpc payload", zap.Error(err))
			return
		}
		c := *p.mqttClient
		if token := c.Publish(p.metricsTopic, 1, false, body); token.Wait() && token.Error() != nil {
			p.logger.Error("error sending metrics RPC", zap.String("topic", p.metricsTopic), zap.Error(token.Error()))
			return
		}

		p.logger.Info("scraped and published metrics", zap.String("topic", p.metricsTopic), zap.Int("payload_size_b", totalSize), zap.Int("batch_count", len(batchPayload)))

	})

	if err != nil {
		return err
	}

	job.SingletonMode()
	return nil
}

func (p *pktvisorBackend) receiveOtlp() {
	exeCtx, execCancelF := context.WithCancel(p.ctx)
	go func() {
		defer execCancelF()
		var err error
		count := 0
		for {
			if p.mqttClient != nil {
				p.exporter, err = p.createOtlpMqttExporter(exeCtx, execCancelF)
				if err != nil {
					p.logger.Error("failed to create a exporter", zap.Error(err))
					return
				}
				pFactory := otlpreceiver.NewFactory()
				cfg := pFactory.CreateDefaultConfig()
				cfg.(*otlpreceiver.Config).Protocols = otlpreceiver.Protocols{
					HTTP: &confighttp.HTTPServerSettings{
						Endpoint: p.otelReceiverHost + ":" + strconv.Itoa(p.otelReceiverPort),
					},
				}
				set := receiver.CreateSettings{
					TelemetrySettings: component.TelemetrySettings{
						Logger:         p.logger,
						TracerProvider: trace.NewNoopTracerProvider(),
						MeterProvider:  global.MeterProvider(),
					},
					BuildInfo: component.NewDefaultBuildInfo(),
				}
				p.receiver, err = pFactory.CreateMetricsReceiver(exeCtx, set, cfg, p.exporter)
				if err != nil {
					p.logger.Error("failed to create a receiver", zap.Error(err))
					return
				}
				err = p.exporter.Start(exeCtx, nil)
				if err != nil {
					p.logger.Error("otel mqtt exporter startup error", zap.Error(err))
					return
				}
				err = p.receiver.Start(exeCtx, nil)
				if err != nil {
					p.logger.Error("otel receiver startup error", zap.Error(err))
					return
				}
				break
			} else {
				count++
				p.logger.Info("waiting until mqtt client is connected try " + strconv.Itoa(count) + " from 10")
				time.Sleep(time.Second * 3)
				if count >= 10 {
					execCancelF()
					_ = p.Stop(exeCtx)
					break
				}
			}
		}
		for {
			select {
			case <-exeCtx.Done():
				p.ctx.Done()
				p.cancelFunc()
			case <-p.ctx.Done():
				err := p.exporter.Shutdown(exeCtx)
				if err != nil {
					return
				}
				err = p.receiver.Shutdown(exeCtx)
				if err != nil {
					return
				}
				p.logger.Info("stopped Orb OpenTelemetry agent collector")
				return
			}
		}
	}()
}
