/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package otelinf

import "github.com/orb-community/orb/pkg/types"

const CurrentSchemaVersion = "1.0"

// policy for otelinf
type collectionPolicy struct {
	Config types.Metadata `json:"config"`
	Kind   string         `json:"kind"`
}

/*

config:
  receivers:
    prometheus:
      config:
        scrape_configs:
          - job_name: 'otelcollector'
            scrape_interval: 1m
            static_configs:
              - targets: ['127.0.0.1:8888']
  exporters:
    otlp:
      endpoint: 127.0.0.1:4317
      tls:
        insecure: true
  service:
    pipelines:
      metrics:
        receivers:
          - prometheus
        exporters:
          - otlp
kind: collection

*/
