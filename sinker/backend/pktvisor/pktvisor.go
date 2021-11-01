/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/pb"
	"github.com/ns1labs/orb/sinker/backend"
	"github.com/ns1labs/orb/sinker/prometheus"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	logger *zap.Logger
}

type context struct {
	agent    *pb.OwnerRes
	agentID  string
	policyID string
}

func (p pktvisorBackend) ProcessMetrics(agent *pb.OwnerRes, agentID string, channelID string, subtopic []string, payload []fleet.AgentMetricsRPCPayload) ([]prometheus.TimeSeries, error) {
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
		context := context{
			agent:    agent,
			agentID:  agentID,
			policyID: data.PolicyID,
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
		tsList = append(tsList, parseToProm(&context, stats)...)
	}
	return tsList, nil
}

func parseToProm(ctxt *context, stats StatSnapshot) prometheus.TSList {
	var tsList = prometheus.TSList{}
	statsMap := structs.Map(stats)
	convertToPromParticle(ctxt, statsMap, "", &tsList)
	return tsList
}

func convertToPromParticle(ctxt *context, m map[string]interface{}, label string, tsList *prometheus.TSList) {
	for k, v := range m {
		switch c := v.(type) {
		case map[string]interface{}:
			convertToPromParticle(ctxt, c, label+k, tsList)
		case int64:
			{
				var matchFirstQuantile = regexp.MustCompile("^([P-p])+[0-9]")
				if ok := matchFirstQuantile.MatchString(k); ok {
					tsList = makePromParticle(ctxt, label, k, v, tsList, ok)
				} else {
					tsList = makePromParticle(ctxt, label+k, "", v, tsList, false)
				}
			}
		case []interface{}:
			{
				for _, value := range c {
					m, ok := value.(map[string]interface{})
					if !ok {
						return
					}
					var lbl string
					var dtpt interface{}
					for k, v := range m {
						switch k {
						case "Name":
							{
								lbl = fmt.Sprintf("%v", v)
							}
						case "Estimate":
							{
								dtpt = v
							}
						}
					}
					tsList = makePromParticle(ctxt, label+k, lbl, dtpt, tsList, false)
				}
			}
		}
	}
}

func makePromParticle(ctxt *context, label string, k string, v interface{}, tsList *prometheus.TSList, quantile bool) *prometheus.TSList {
	mapQuantiles := make(map[string]float64)
	mapQuantiles["P50"] = 0.50
	mapQuantiles["P90"] = 0.90
	mapQuantiles["P95"] = 0.95
	mapQuantiles["P99"] = 0.99

	var dpFlag dp
	var labelsListFlag labelList
	labelsListFlag.Set(fmt.Sprintf("__name__:%s", camelToSnake(label)))
	labelsListFlag.Set("instance:" + ctxt.agent.AgentName)
	labelsListFlag.Set("agent_id:" + ctxt.agentID)
	labelsListFlag.Set("policy_id:" + ctxt.policyID)
	if k != "" {
		if quantile {
			if value, ok := mapQuantiles[k]; ok {
				labelsListFlag.Set(fmt.Sprintf("quantile:%.2f", value))
			}
		} else {
			labelsListFlag.Set(fmt.Sprintf("name:%s", k))
		}
	}
	dpFlag.Set(fmt.Sprintf("now,%d", v))
	*tsList = append(*tsList, prometheus.TimeSeries{
		Labels:    []prometheus.Label(labelsListFlag),
		Datapoint: prometheus.Datapoint(dpFlag),
	})
	return tsList
}

func camelToSnake(s string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	lower := strings.ToLower(snake)
	return lower
}

func Register(logger *zap.Logger) bool {
	backend.Register("pktvisor", &pktvisorBackend{logger: logger})
	return true
}
