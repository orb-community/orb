package cloudprober

import (
	"encoding/json"
	"github.com/ns1labs/orb/agent/policies"
	"go.uber.org/zap"
	"testing"
)

func Test_cloudproberBackend_buildConfigFile(t *testing.T) {
	t.Skip("local only")
	type args struct {
		policyYaml []policies.PolicyData
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "single policy",
			args: args{
				policyYaml: []policies.PolicyData{
					{
						ID:      "singlepolicy-id",
						Name:    "singlepolicy-name",
						Backend: BackendName,
						Data: map[string]interface{}{
							"probes": []map[string]interface{}{
								{
									"name":               "simple_google_probe",
									"type":               "http",
									"interval_msec":      2000,
									"timeout_msec":       3000,
									"targets_host_names": "www.google.com,ns1.com,www.reddit.com",
								},
							},
						},
					},
				},
			},
			want:    []byte("probe {\n  name: \"sample_google_probe\"\n  type: HTTP\n  targets {\n    host_names: \"www.google.com\"\n  }\n  interval_msec: 5000  # 5s\n  timeout_msec: 1000   # 1s\n}\nserver {\n  type: HTTP\n  http_server {\n    port: 8099\n  }\n}\n\nsurfacer {\n  type: PROMETHEUS\n\n  prometheus_surfacer {\n    # Following option adds a prefix to exported metrics, for example,\n    # \"total\" metric is exported as \"cloudprober_total\".\n    metrics_prefix: \"cloudprober_\"\n  }\n}"),
			wantErr: false,
		},
		{
			name: "multiple probe policy",
			args: args{
				policyYaml: []policies.PolicyData{
					{
						ID:      "multipleprobe-id",
						Name:    "multipleprobe-name",
						Backend: BackendName,
						Data: map[string]interface{}{
							"probes": []map[string]interface{}{
								{
									"name":               "multi_google_probe",
									"type":               "http",
									"interval_msec":      5000,
									"timeout_msec":       1000,
									"targets_host_names": "www.google.com",
								},
								{
									"name":               "multi_ping_probe",
									"type":               "ping",
									"interval_msec":      5000,
									"timeout_msec":       1000,
									"targets_host_names": "www.google.com,ns1.com,www.reddit.com",
								},
							},
						},
					},
				},
			},
			want:    []byte("probe {\n  name: \"sample_google_probe\"\n  type: HTTP\n  targets {\n    host_names: \"www.google.com\"\n  }\n  interval_msec: 5000  # 5s\n  timeout_msec: 1000   # 1s\n}\nserver {\n  type: HTTP\n  http_server {\n    port: 8099\n  }\n}\n\nsurfacer {\n  type: PROMETHEUS\n\n  prometheus_surfacer {\n    # Following option adds a prefix to exported metrics, for example,\n    # \"total\" metric is exported as \"cloudprober_total\".\n    metrics_prefix: \"cloudprober_\"\n  }\n}"),
			wantErr: false,
		},
		{
			name: "parse multiple policies",
			args: args{
				policyYaml: []policies.PolicyData{
					{
						ID:      "singlepolicy-id",
						Name:    "singlepolicy-name",
						Backend: BackendName,
						Data: map[string]interface{}{
							"probes": []map[string]interface{}{
								{
									"name":               "simple_google_probe",
									"type":               "http",
									"interval_msec":      5000,
									"timeout_msec":       1000,
									"targets_host_names": "www.google.com",
								},
							},
						},
					},
					{
						ID:      "multipleprobe-id",
						Name:    "multipleprobe-name",
						Backend: BackendName,
						Data: map[string]interface{}{
							"probes": []map[string]interface{}{
								{
									"name":               "multi_google_probe",
									"type":               "http",
									"interval_msec":      5000,
									"timeout_msec":       1000,
									"targets_host_names": "www.google.com",
								},
								{
									"name":               "multi_ping_probe",
									"type":               "ping",
									"interval_msec":      5000,
									"timeout_msec":       1000,
									"targets_host_names": "www.google.com,ns1.com,www.reddit.com",
								},
							},
						},
					},
				},
			},
			want:    []byte("probe {\n  name: \"sample_google_probe\"\n  type: HTTP\n  targets {\n    host_names: \"www.google.com\"\n  }\n  interval_msec: 5000  # 5s\n  timeout_msec: 1000   # 1s\n}\nserver {\n  type: HTTP\n  http_server {\n    port: 8099\n  }\n}\n\nsurfacer {\n  type: PROMETHEUS\n\n  prometheus_surfacer {\n    # Following option adds a prefix to exported metrics, for example,\n    # \"total\" metric is exported as \"cloudprober_total\".\n    metrics_prefix: \"cloudprober_\"\n  }\n}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cloudproberBackend{
				logger:     zap.NewNop(),
				configFile: "",
			}
			got, err := c.buildConfigFile(tt.args.policyYaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got ConfigFile: %s", got)

		})
	}
}

func Test_cloudproberBackend_parsingYaml(t *testing.T) {
	t.Skip("local only")
	type args struct {
		policyJson string
	}
	tests := []struct {
		name    string
		args    args
		want    Probes
		wantErr bool
	}{
		{
			name: "testing parse single policy",
			args: args{
				policyJson: `{
	"probes": {
		"p_name": {
			"type": "HTTP",
			"interval_msec": 2000,
			"timeout_msec": 3000,
			"targets_host_names": "www.google.com,ns1.com"
		 }
	}
}`,
			},
			want: Probes{ProbeData: []ProbeData{
				{
					Name:         "",
					ProbeType:    "HTTP",
					Targets:      "www.google.com,ns1.com",
					IntervalMsec: 2000,
					TimeoutMsec:  3000,
				},
			},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *Probes
			err := json.Unmarshal([]byte(tt.args.policyJson), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got ConfigFile: %s", got.ToConfigFile())
		})
	}
}
