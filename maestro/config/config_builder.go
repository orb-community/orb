package config

import (
	"context"
	"fmt"
	"github.com/ns1labs/orb/pkg/errors"
	"gopkg.in/yaml.v2"
	"strings"
)

var k8sOtelCollector = `
{
  "kind": "List",
  "apiVersion": "v1",
  "metadata": {
  },
  "items": [
    {
      "kind": "ConfigMap",
      "apiVersion": "v1",
      "metadata": {
        "name": "otel-collector-config-SINK_ID",
        "creationTimestamp": null
      },
      "data": {
        "config.yaml": "SINK_CONFIG"
      }
    },
    {
      "kind": "Deployment",
      "apiVersion": "apps/v1",
      "metadata": {
        "name": "otel-SINK_ID",
        "creationTimestamp": null,
        "labels": {
          "app": "opentelemetry",
          "component": "otel-collector"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "app": "opentelemetry",
            "component": "otel-collector-SINK_ID"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "app": "opentelemetry",
              "component": "otel-collector-SINK_ID"
            }
          },
          "spec": {
            "volumes": [
              {
                "name": "varlog",
                "hostPath": {
                  "path": "/var/log",
                  "type": ""
                }
              },
              {
                "name": "varlibdockercontainers",
                "hostPath": {
                  "path": "/var/lib/docker/containers",
                  "type": ""
                }
              },
              {
                "name": "data",
                "configMap": {
                  "name": "otel-collector-config-SINK_ID",
                  "defaultMode": 420
                }
              }
            ],
            "containers": [
              {
                "name": "otel-collector",
                "image": "otel/opentelemetry-collector-contrib:0.68.0",
                "ports": [
                  {
                    "containerPort": 13133,
                    "protocol": "TCP"
                  },
                  {
                    "containerPort": 8888,
                    "protocol": "TCP"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "100m",
                    "memory": "200Mi"
                  },
                  "requests": {
                    "cpu": "100m",
                    "memory": "200Mi"
                  }
                },
                "volumeMounts": [
                  {
                    "name": "varlog",
                    "readOnly": true,
                    "mountPath": "/var/log"
                  },
                  {
                    "name": "varlibdockercontainers",
                    "readOnly": true,
                    "mountPath": "/var/lib/docker/containers"
                  },
                  {
                    "name": "data",
                    "readOnly": true,
                    "mountPath": "/etc/otelcol-contrib/config.yaml",
                    "subPath": "config.yaml"
                  }
                ],
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "imagePullPolicy": "IfNotPresent"
              }
            ],
            "restartPolicy": "Always",
            "terminationGracePeriodSeconds": 30,
            "dnsPolicy": "ClusterFirst",
            "securityContext": {
            },
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": "25%",
            "maxSurge": "25%"
          }
        },
        "revisionHistoryLimit": 10,
        "progressDeadlineSeconds": 600
      },
      "status": {
      }
    },
    {
      "kind": "Service",
      "apiVersion": "v1",
      "metadata": {
        "name": "otel-SINK_ID",
        "creationTimestamp": null,
        "labels": {
          "app": "opentelemetry",
          "component": "otel-collector-SINK_ID"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "metrics",
            "protocol": "TCP",
            "port": 8888,
            "targetPort": 8888
          },
          {
            "name": "healthcheck",
            "protocol": "TCP",
            "port": 13133,
            "targetPort": 13133
          }
        ],
        "selector": {
          "component": "otel-collector-SINK_ID"
        },
        "type": "ClusterIP",
        "sessionAffinity": "None"
      },
      "status": {
        "loadBalancer": {
        }
      }
    }
  ]
}
`

func GetDeploymentJson(kafkaUrl, sinkId, sinkUrl, sinkUsername, sinkPassword string) (string, error) {
	// prepare manifest
	manifest := strings.Replace(k8sOtelCollector, "SINK_ID", sinkId, -1)
	config, err := ReturnConfigYamlFromSink(context.Background(), kafkaUrl, sinkId, sinkUrl, sinkUsername, sinkPassword)
	if err != nil {
		return "", errors.Wrap(errors.New("failed to build YAML"), err)
	}
	manifest = strings.Replace(manifest, "SINK_CONFIG", config, -1)
	return manifest, nil
}

// ReturnConfigYamlFromSink this is the main method, which will generate the YAML file from the
func ReturnConfigYamlFromSink(_ context.Context, kafkaUrlConfig, sinkId, sinkUrl, sinkUsername, sinkPassword string) (string, error) {
	config := OtelConfigFile{
		Receivers: Receivers{
			Kafka: KafkaReceiver{
				Brokers:         []string{kafkaUrlConfig},
				Topic:           fmt.Sprintf("otlp_metrics-%s", sinkId),
				ProtocolVersion: "2.0.0", // Leaving default of over 2.0.0
			},
		},
		Extensions: &Extensions{
			HealthCheckExtConfig: &HealthCheckExtension{
				Endpoint: "0.0.0.0:13133",
				Path:     "/health/status",
				CollectorPipeline: &CollectorPipelineExtension{
					Enabled: "true", Interval: "5m", FailureThreshold: 4},
			},
			PProf: &PProfExtension{
				Endpoint: "0.0.0.0:1888", // Leaving default for now, will need to change with more processes
			},
			BasicAuth: &BasicAuthenticationExtension{
				ClientAuth: &struct {
					Username string `json:"username" yaml:"username"`
					Password string `json:"password" yaml:"password"`
				}{Username: sinkUsername, Password: sinkPassword},
			},
		},
		Exporters: Exporters{
			PrometheusRemoteWrite: &PrometheusRemoteWriteExporterConfig{
				Endpoint: sinkUrl,
				Auth: struct {
					Authenticator string `json:"authenticator" yaml:"authenticator"`
				}{Authenticator: "basicauth/exporter"},
			},
		},
		Service: ServiceConfig{
			Extensions: []string{"pprof", "health_check", "basicauth/exporter"},
			Pipelines: struct {
				Metrics struct {
					Receivers  []string `json:"receivers" yaml:"receivers"`
					Processors []string `json:"processors,omitempty" yaml:"processors,omitempty"`
					Exporters  []string `json:"exporters" yaml:"exporters"`
				} `json:"metrics" yaml:"metrics"`
			}{
				Metrics: struct {
					Receivers  []string `json:"receivers" yaml:"receivers"`
					Processors []string `json:"processors,omitempty" yaml:"processors,omitempty"`
					Exporters  []string `json:"exporters" yaml:"exporters"`
				}{
					Receivers: []string{"kafka"},
					Exporters: []string{"prometheusremotewrite"},
				},
			},
		},
	}
	marshal, err := yaml.Marshal(&config)
	if err != nil {
		return "", err
	}
	returnedString := "---\n" + string(marshal)
	s := strings.ReplaceAll(returnedString, "\"", "")
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s, nil

}

type OtelConfigFile struct {
	Receivers  Receivers     `json:"receivers" yaml:"receivers"`
	Processors *Processors   `json:"processors,omitempty" yaml:"processors,omitempty"`
	Extensions *Extensions   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Exporters  Exporters     `json:"exporters" yaml:"exporters"`
	Service    ServiceConfig `json:"service" yaml:"service"`
}

// Receivers will receive only with Kafka for now
type Receivers struct {
	Kafka KafkaReceiver `json:"kafka" yaml:"kafka"`
}

type KafkaReceiver struct {
	Brokers         []string `json:"brokers" yaml:"brokers"`
	Topic           string   `json:"topic" yaml:"topic"`
	ProtocolVersion string   `json:"protocol_version" yaml:"protocol_version"`
}

type Processors struct {
}

type Extensions struct {
	HealthCheckExtConfig *HealthCheckExtension `json:"health_check,omitempty" yaml:"health_check,omitempty"`
	PProf                *PProfExtension       `json:"pprof,omitempty" yaml:"pprof,omitempty"`
	ZPages               *ZPagesExtension      `json:"zpages,omitempty" yaml:"zpages,omitempty"`
	// Exporters Authentication
	BasicAuth *BasicAuthenticationExtension `json:"basicauth/exporter,omitempty" yaml:"basicauth/exporter,omitempty"`
}

type HealthCheckExtension struct {
	Endpoint          string                      `json:"endpoint" yaml:"endpoint"`
	Path              string                      `json:"path" yaml:"path"`
	CollectorPipeline *CollectorPipelineExtension `json:"check_collector_pipeline,omitempty" yaml:"check_collector_pipeline,omitempty"`
}

type CollectorPipelineExtension struct {
	Enabled          string `json:"enabled" yaml:"enabled"`
	Interval         string `json:"interval" yaml:"interval"`
	FailureThreshold int32  `json:"exporter_failure_threshold" yaml:"exporter_failure_threshold"`
}

type PProfExtension struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type ZPagesExtension struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type BasicAuthenticationExtension struct {
	ClientAuth *struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	} `json:"client_auth" yaml:"client_auth"`
}

type Exporters struct {
	PrometheusRemoteWrite *PrometheusRemoteWriteExporterConfig `json:"prometheusremotewrite,omitempty" yaml:"prometheusremotewrite,omitempty"`
}

type PrometheusRemoteWriteExporterConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Auth     struct {
		Authenticator string `json:"authenticator" yaml:"authenticator"`
	}
}

type ServiceConfig struct {
	Extensions []string `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	Pipelines  struct {
		Metrics struct {
			Receivers  []string `json:"receivers" yaml:"receivers"`
			Processors []string `json:"processors,omitempty" yaml:"processors,omitempty"`
			Exporters  []string `json:"exporters" yaml:"exporters"`
		} `json:"metrics" yaml:"metrics"`
	} `json:"pipelines" yaml:"pipelines"`
}
