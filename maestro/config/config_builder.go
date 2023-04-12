package config

import (
	"context"
	"fmt"
	"github.com/orb-community/orb/pkg/errors"
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

var JsonService = `
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
`

var JsonConfigMap = `
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
}
`

var JsonDeployment = `
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

func GetDeploymentApplyConfig(sinkId string) string {
	manifest := strings.Replace(JsonDeployment, "SINK_ID", sinkId, -1)
	return manifest
}

func GetConfigMapApplyConfig(kafkaUrl, sinkId, sinkUrl, sinkUsername, sinkPassword string) (string, error) {
	manifest := strings.Replace(JsonConfigMap, "SINK_ID", sinkId, -1)
	config, err := ReturnConfigYamlFromSink(context.Background(), kafkaUrl, sinkId, sinkUrl, sinkUsername, sinkPassword)
	if err != nil {
		return "", errors.Wrap(errors.New("failed to build YAML"), err)
	}
	manifest = strings.Replace(manifest, "SINK_CONFIG", config, -1)
	return manifest, nil
}

func GetServiceApplyConfig(sinkId string) string {
	manifest := strings.Replace(JsonService, "SINK_ID", sinkId, -1)
	return manifest
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
			LoggingExporter: &LoggingExporterConfig{
				Verbosity:          "detailed",
				SamplingInitial:    5,
				SamplingThereAfter: 50,
			},
		},
		Service: ServiceConfig{
			Extensions: []string{"pprof", "basicauth/exporter"},
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
