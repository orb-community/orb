package otlpmqttexporter_test

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/orb-community/orb/agent/otel/otlpmqttexporter"
	"github.com/orb-community/orb/agent/policies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractCollectorConfig(t *testing.T) {
	policyJSON := `{
  "ID": "test-policy-id",
  "Datasets": null,
  "GroupIds": null,
  "Name": "test-policy",
  "Backend": "otel",
  "Version": 0,
  "Format": "yaml",
  "Data": {
    "exporters": {},
    "receivers": {
      "httpcheck": {
        "collection_interval": "60s",
        "targets": [
          {
            "endpoint": "https://example.com",
            "method": "GET",
            "tags": {
              "foo": "bar"
            }
          }
        ]
      }
    },
    "service": {
      "pipelines": {
        "metrics": {
          "exporters": null,
          "receivers": [
            "httpcheck"
          ]
        }
      }
    }
  },
  "State": 1,
  "BackendErr": "",
  "LastScrapeBytes": 0,
  "LastScrapeTS": "2023-12-18T13:57:42.024296Z",
  "PreviousPolicyData": null
}`
	var policy policies.PolicyData
	if err := json.Unmarshal([]byte(policyJSON), &policy); err != nil {
		t.Fatal(err)
	}

	cfg, err := otlpmqttexporter.ExtractCollectorConfig(policy)
	require.NoError(t, err)
	assert.Equal(t, 1, len(cfg.Receivers))

	for key, value := range cfg.Receivers {
		switch key.Type() {
		case "httpcheck":
			var httpcheck otlpmqttexporter.HTTPCheckReceiver
			if err := mapstructure.Decode(value, &httpcheck); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, "60s", httpcheck.CollectionInterval)
			assert.Equal(t, 1, len(httpcheck.Targets))
			for _, target := range httpcheck.Targets {
				assert.Equal(t, "https://example.com", target.Endpoint)
				assert.Equal(t, "GET", target.Method)
				assert.Equal(t, map[string]string{"foo": "bar"}, target.Tags)
			}
		}
	}
}
