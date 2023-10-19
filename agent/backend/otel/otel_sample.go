package otel

import (
	"github.com/orb-community/orb/agent/policies"
	"time"
)

var samplePolicyData = `---
receivers:
  httpcheck:
	targets:
     - endpoint: http://localhost:8000/health
     - method: GET
exporters:
  otlphttp:
    endpoint: http://localhost:0
service:
  pipelines:
    metrics:
      exporters: 
        - otlphttp
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
