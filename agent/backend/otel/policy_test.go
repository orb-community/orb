package otel_test

import (
	"testing"
	"time"

	"github.com/orb-community/orb/agent/backend/otel"
	"github.com/orb-community/orb/agent/policies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestSanitizePolicyData(t *testing.T) {
	newPolicyData := policies.PolicyData{
		ID:              "test-policy-id",
		Name:            "test-policy",
		Backend:         "otel",
		Version:         0,
		Format:          "yaml",
		State:           policies.Running,
		LastScrapeBytes: 0,
		LastScrapeTS:    time.Now(),
		Data: map[string]interface{}{
			"receivers": map[string]interface{}{
				"httpcheck": map[string]interface{}{
					"collection_interval": "60s",
					"targets": []interface{}{
						map[string]interface{}{
							"endpoint": "https://example.com",
							"method":   "GET",
							"tags": map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			"exporters": map[string]interface{}{},
			"service": map[string]interface{}{
				"pipelines": map[string]interface{}{
					"metrics": map[string]interface{}{
						"exporters": nil,
						"receivers": []string{"httpcheck"},
					},
				},
			},
		},
		PreviousPolicyData: nil,
	}

	policyYaml, err := yaml.Marshal(newPolicyData.Data)
	require.NoError(t, err)

	copyPolicyData, err := otel.SanitizePolicyData(newPolicyData)
	require.NoError(t, err)

	copyPolicyYaml, err := yaml.Marshal(copyPolicyData.Data)
	require.NoError(t, err)

	assert.NotEqual(t, newPolicyData, copyPolicyData)
	assert.NotEqual(t, string(policyYaml), string(copyPolicyYaml))
}
