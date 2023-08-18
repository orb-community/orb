package otel

import (
	"github.com/orb-community/orb/agent/policies"
	"time"
)

var samplePolicyData = `---
receivers:
  httpcheck:
    endpoint: http://localhost:8000/health
    method: GET
    collection_interval: 1s
exporters:
  otlphttp:
    endpoint: http://localhost:0
  logging:
    verbosity: detailed
    sampling_initial: 10
    sampling_thereafter: 200
service:
  pipelines:
    metrics:
      exporters: 
        - otlphttp
        - logging
      receivers: 
        - httpcheck
`

var samplePolicy = policies.PolicyData{
	ID:                 "default",
	Datasets:           nil,
	GroupIds:           nil,
	Name:               "opentelemetry-default",
	Backend:            "otel",
	Version:            0,
	Data:               samplePolicyData,
	State:              0,
	BackendErr:         "",
	LastScrapeBytes:    0,
	LastScrapeTS:       time.Time{},
	PreviousPolicyData: nil,
}
