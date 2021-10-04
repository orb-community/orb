/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/sinker/backend"
	"github.com/ns1labs/orb/sinker/prometheus"
	"go.uber.org/zap"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	logger *zap.Logger
}

func (p pktvisorBackend) ProcessMetrics(thingID string, channelID string, subtopic []string, payload []fleet.AgentMetricsRPCPayload) ([]prometheus.TimeSeries, error) {
	// process batch
	var tsList = []prometheus.TimeSeries{}
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
		tsList = append(tsList, parseToProm(stats)...)
	}
	return tsList, nil
}

func parseToProm(stats StatSnapshot) prometheus.TSList {
	//var headerListFlag headerList
	//var headers map[string]string
	//if len(headerListFlag) > 0 {
	//	log.Println("with headers", headerListFlag.String())
	//	headers = make(map[string]string, len(headerListFlag))
	//	for _, header := range headerListFlag {
	//		headers[header.name] = header.value
	//	}
	//}

	var tsList = prometheus.TSList{}
	for _, v := range stats.DNS.TopUDPPorts {
		var dpFlag dp
		var labelsListFlag labelList
		labelsListFlag.Set("__name__:dns_top_udp_ports")
		labelsListFlag.Set("instance:gw")
		labelsListFlag.Set(fmt.Sprintf("name:%s", v.Name))
		dpFlag.Set(fmt.Sprintf("now,%d", v.Estimate))
		tsList = append(tsList, prometheus.TimeSeries{
			Labels:    []prometheus.Label(labelsListFlag),
			Datapoint: prometheus.Datapoint(dpFlag),
		})
	}

	for _, v := range stats.DNS.TopQname2 {
		var dpFlag dp
		var labelsListFlag labelList
		labelsListFlag.Set("__name__:dns_top_qname2")
		labelsListFlag.Set("instance:gw")
		labelsListFlag.Set(fmt.Sprintf("name:%s", v.Name))
		dpFlag.Set(fmt.Sprintf("now,%d", v.Estimate))
		tsList = append(tsList, prometheus.TimeSeries{
			Labels:    []prometheus.Label(labelsListFlag),
			Datapoint: prometheus.Datapoint(dpFlag),
		})
	}

	return tsList
}

func Register(logger *zap.Logger) bool {
	backend.Register("pktvisor", &pktvisorBackend{logger: logger})
	return true
}
