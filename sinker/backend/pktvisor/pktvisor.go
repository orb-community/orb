/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/sinker/backend"
	"go.uber.org/zap"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	logger *zap.Logger
}

func (p pktvisorBackend) ProcessMetrics(thingID string, channelID string, subtopic []string, payload []fleet.AgentMetricsRPCPayload) error {
	// process batch
	for _, data := range payload {
		// TODO check pktvisor version in data.BEVersion against PktvisorVersion
		// TODO policyID and datasetIDs are in data
		if data.Format != "json" {
			p.logger.Warn("ignoring non-json pktvisor payload", zap.String("format", data.Format))
			continue
		}
		// unmarshal pktvisor metrics
		var metrics map[string]map[string]interface{}
		err := json.Unmarshal(data.Data, &metrics)
		if err != nil {
			p.logger.Warn("unable to unmarshal pktvisor metric payload", zap.Any("payload", data.Data))
			continue
		}
		stats := StatSnapshot{}
		for _, handlerData := range metrics {
			if data, ok := handlerData["pcap"]; ok {
				err := mapstructure.Decode(data, &stats.Pcap)
				if err != nil {
					p.logger.Error("error decoding pcap handler", zap.Error(err))
					continue
				}
			} else if data, ok := handlerData["dns"]; ok {
				err := mapstructure.Decode(data, &stats.DNS)
				if err != nil {
					p.logger.Error("error decoding dns handler", zap.Error(err))
					continue
				}
			} else if data, ok := handlerData["packets"]; ok {
				err := mapstructure.Decode(data, &stats.Packets)
				if err != nil {
					p.logger.Error("error decoding packets handler", zap.Error(err))
					continue
				}
			}
		}
		// TODO turn StatSnapshot into format for prometheus remote_write
		p.logger.Info("decoded pktvisor metrics", zap.Any("metrics", stats))
	}
	return nil
}

func Register(logger *zap.Logger) bool {
	backend.Register("pktvisor", &pktvisorBackend{logger: logger})
	return true
}
