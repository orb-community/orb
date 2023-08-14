package otel

import (
	"github.com/orb-community/orb/agent/policies"
	"time"
)

var policyData = `
receivers:
  httpcheck:
    targets:
      - endpoint: http://localhost:8000/health
        method: GET
      - endpoint: http://localhost:8000/health
        method: GET
    collection_interval: 10s

exporters:
	otlphttp:
		endpoint: http://localhost:

processors: # this collection value may be supported (tbd)

extensions: # tbd

service: # tbd
	metrics:
		exporters: 
`

var policy = policies.PolicyData{
	ID:                 "default",
	Datasets:           nil,
	GroupIds:           nil,
	Name:               "opentelemetry-default",
	Backend:            "otel",
	Version:            0,
	Data:               nil,
	State:              0,
	BackendErr:         "",
	LastScrapeBytes:    0,
	LastScrapeTS:       time.Time{},
	PreviousPolicyData: nil,
}
