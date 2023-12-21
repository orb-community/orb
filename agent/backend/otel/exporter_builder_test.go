package otel

import (
	"go.uber.org/zap"
	"testing"
)

func TestBuildDefaultPolicy(t *testing.T) {
	testCases := []struct {
		caseName        string
		inputString     string
		policyId        string
		policyName      string
		expectedStruct  openTelemetryConfig
		processedString string
		wantErr         error
	}{
		{
			caseName: "success default policy test",
			inputString: `
---
receivers:
  httpcheck:
    targets:
      - endpoint: http://orb.live
        method: GET
      - endpoint: http://orb.community
        method: GET
    collection_interval: 60s
exporters:
  otlp:
    endpoint: localhost:4317
    tls:
      insecure: true
  logging:
    verbosity: detailed
    sampling_initial: 5
service:
  pipelines:
    metrics:
      exporters: 
        - otlp
      receivers: 
        - httpcheck
`,
			policyId:   "test-policy-id",
			policyName: "test-policy",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			logger := zap.NewNop()
			exporterBuilder := getExporterBuilder(logger, "localhost", 4317)
			gotOtelConfig, err := exporterBuilder.GetStructFromYaml(testCase.inputString)
			if err != nil {
				t.Errorf("failed to merge default value with policy: %v", err)
			}
			expectedStruct, err := exporterBuilder.MergeDefaultValueWithPolicy(gotOtelConfig, testCase.policyId, testCase.policyName)
			if err != nil {
				t.Errorf("failed to merge default value with policy: %v", err)
			}
			if _, ok := expectedStruct.Processors["transform/policy_data"]; !ok {
				t.Error("missing required attributes/policy_data processor", err)
			}

		})
	}
}
