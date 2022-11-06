/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package backend

import (
	"github.com/etaques/orb/fleet"
	"github.com/etaques/orb/fleet/pb"
	"github.com/etaques/orb/sinker/prometheus"
)

type Backend interface {
	ProcessMetrics(agent *pb.AgentInfoRes, thingID string, data fleet.AgentMetricsRPCPayload) ([]prometheus.TimeSeries, error)
}

var registry = make(map[string]Backend)

func Register(name string, b Backend) {
	registry[name] = b
}

func GetList() []string {
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}
	return keys
}

func HaveBackend(name string) bool {
	_, prs := registry[name]
	return prs
}

func GetBackend(name string) Backend {
	return registry[name]
}
