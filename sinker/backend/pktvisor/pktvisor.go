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
	agent      *pb.OwnerRes
	agentID    string
	policyID   string
	policyName string
	logger     *zap.Logger
}

func (p pktvisorBackend) ProcessMetrics(agent *pb.OwnerRes, agentID string, data fleet.AgentMetricsRPCPayload) ([]prometheus.TimeSeries, error) {
	// TODO check pktvisor version in data.BEVersion against PktvisorVersion
	if data.Format != "json" {
		p.logger.Warn("ignoring non-json pktvisor payload", zap.String("format", data.Format))
		return nil, nil
	}
	// unmarshal pktvisor metrics
	var metrics map[string]map[string]interface{}
	err := json.Unmarshal(data.Data, &metrics)
	if err != nil {
		p.logger.Warn("unable to unmarshal pktvisor metric payload", zap.Any("payload", data.Data))
		return nil, err
	}
	context := context{
		agent:      agent,
		agentID:    agentID,
		policyID:   data.PolicyID,
		policyName: data.PolicyName,
		logger:     p.logger,
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
	return parseToProm(&context, stats), nil
}

func parseToProm(ctxt *context, stats StatSnapshot) prometheus.TSList {
	var tsList = prometheus.TSList{}
	statsMap := structs.Map(stats)
	convertToPromParticle(ctxt, statsMap, "", &tsList)
	return tsList
}

func convertToPromParticle(ctxt *context, statsMap map[string]interface{}, label string, tsList *prometheus.TSList) {
	for key, value := range statsMap {
		switch statistic := value.(type) {
		case map[string]interface{}:
			// Call convertToPromParticle recursively until the last interface of the StatSnapshot struct
			// The prom particle label it's been formed during the recursive call (concatenation)
			convertToPromParticle(ctxt, statistic, label+key, tsList)
		// The StatSnapshot has two ways to record metrics (i.e. Live int64 `mapstructure:"live"`)
		// It's why we check if the type is int64
		case int64:
			{
				// Use this regex to identify if the value it's a quantile
				var matchFirstQuantile = regexp.MustCompile("^([P-p])+[0-9]")
				if ok := matchFirstQuantile.MatchString(key); ok {
					// If it's quantile, needs to be parsed to prom quantile format
					tsList = makePromParticle(ctxt, label, key, value, tsList, ok, "")
				} else {
					tsList = makePromParticle(ctxt, label+key, "", value, tsList, false, "")
				}
			}
		// The StatSnapshot has two ways to record metrics (i.e. TopIpv4   []NameCount   `mapstructure:"top_ipv4"`)
		// It's why we check if the type is []interface
		// Here we extract the value for Name and Estimate
		case []interface{}:
			{
				for _, value := range statistic {
					m, ok := value.(map[string]interface{})
					if !ok {
						return
					}
					var promLabel string
					var promDataPoint interface{}
					for k, v := range m {
						switch k {
						case "Name":
							{
								promLabel = fmt.Sprintf("%v", v)
							}
						case "Estimate":
							{
								promDataPoint = v
							}
						}
					}
					tsList = makePromParticle(ctxt, label, promLabel, promDataPoint, tsList, false, key)
				}
			}
		}
	}
}

func makePromParticle(ctxt *context, label string, k string, v interface{}, tsList *prometheus.TSList, quantile bool, name string) *prometheus.TSList {
	mapQuantiles := make(map[string]float64)
	mapQuantiles["P50"] = 0.5
	mapQuantiles["P90"] = 0.9
	mapQuantiles["P95"] = 0.95
	mapQuantiles["P99"] = 0.99

	var dpFlag dp
	var labelsListFlag labelList
	if err := labelsListFlag.Set(fmt.Sprintf("__name__;%s", camelToSnake(label))); err != nil {
		handleParticleError(ctxt, err)
	}
	if err := labelsListFlag.Set("instance;" + ctxt.agent.AgentName); err != nil {
		handleParticleError(ctxt, err)
	}
	if err := labelsListFlag.Set("agent_id;" + ctxt.agentID); err != nil {
		handleParticleError(ctxt, err)
	}
	if err := labelsListFlag.Set("agent;" + ctxt.agent.AgentName); err != nil {
		handleParticleError(ctxt, err)
	}
	if err := labelsListFlag.Set("policy_id;" + ctxt.policyID); err != nil {
		handleParticleError(ctxt, err)
	}
	if err := labelsListFlag.Set("policy;" + ctxt.policyName); err != nil {
		handleParticleError(ctxt, err)
	}

	if k != "" {
		if quantile {
			if value, ok := mapQuantiles[k]; ok {
				if err := labelsListFlag.Set(fmt.Sprintf("quantile;%.2f", value)); err != nil {
					handleParticleError(ctxt, err)
				}
			}
		} else {
			if err := labelsListFlag.Set(fmt.Sprintf("%s;%s", name, k)); err != nil {
				handleParticleError(ctxt, err)
			}
		}
	}
	if err := dpFlag.Set(fmt.Sprintf("now,%d", v)); err != nil {
		handleParticleError(ctxt, err)
	}
	*tsList = append(*tsList, prometheus.TimeSeries{
		Labels:    labelsListFlag,
		Datapoint: prometheus.Datapoint(dpFlag),
	})
	return tsList
}

func handleParticleError(ctxt *context, err error) {
	ctxt.logger.Error("failed to set prometheus element", zap.Error(err))
}

func camelToSnake(s string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	// Approach to avoid change the values to TopGeoLoc and TopASN
	// Should continue camel case or upper case
	var matchExcept = regexp.MustCompile(`(oLoc$|pASN$)`)
	sub := matchExcept.Split(s, 2)
	var strExcept = ""
	if len(sub) > 1 {
		strExcept = matchExcept.FindAllString(s, 1)[0]
		s = sub[0]
	}

	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	lower := strings.ToLower(snake)
	return lower + strExcept
}

func Register(logger *zap.Logger) bool {
	backend.Register("pktvisor", &pktvisorBackend{logger: logger})
	return true
}
