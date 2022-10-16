// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package maestro

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ns1labs/orb/pkg/errors"
	"go.uber.org/zap"
)

var (
	otelCollectorCfg = `---
receivers:
  kafka:
    brokers:
    - orb-live-stg-kafka.orb-live.svc.cluster.local:9092
    topic: otlp_metrics-sink-id-222
    protocol_version: 2.0.0
extensions:
  health_check:
    endpoint: :13133
    path: /health/status
    check_collector_pipeline:
      enabled: true
      interval: 5m
      exporter_failure_threshold: 5
  basicauth/exporter:
    client_auth:
      username: admin
      password: amanda.joaquina
exporters:
  prometheusremotewrite:
    endpoint: https://prometheus.qa.orb.live/api/v1/write
    auth:
      authenticator: basicauth/exporter
service:
  extensions:
  - health_check
  - basicauth/exporter
  pipelines:
    metrics:
      receivers:
      - kafka
      exporters:
      - prometheusremotewrite
`

	k8sOtelCollector = `
{
    "kind": "List",
    "apiVersion": "v1",
    "metadata": {},
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
                                "image": "otel/opentelemetry-collector-contrib:0.60.0",
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
                        "securityContext": {},
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
            "status": {}
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
                    }
                ],
                "selector": {
                    "component": "otel-collector-SINK_ID"
                },
                "type": "ClusterIP",
                "sessionAffinity": "None"
            },
            "status": {
                "loadBalancer": {}
            }
        }
    ]
}
`
	ErrCreateMaestro   = errors.New("failed to create Otel Collector")
	ErrConflictMaestro = errors.New("Otel collector already exists")
)

func (svc maestroService) collectorDeploy(operation string, namespace string, sinkId string, manifest string, config string) error {
	// prepare manifest
	manifest = strings.Replace(manifest, "SINK_ID", sinkId, -1)
	config = strings.Replace(config, "\n", `\n`, -1)
	manifest = strings.Replace(manifest, "SINK_CONFIG", config, -1)
	fileContent := []byte(manifest)
	err := os.WriteFile("/tmp/otel-collector-"+sinkId+".json", fileContent, 0644)
	if err != nil {
		panic(err)
	}

	// execute action
	cmd := exec.Command("kubectl", operation, "-f", "/tmp/otel-collector-"+sinkId+".json", "-n", namespace)
	stdoutReader, _ := cmd.StdoutPipe()
	stdoutScanner := bufio.NewScanner(stdoutReader)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Println(stdoutScanner.Text())
			svc.logger.Info("Deploy Info: " + stdoutScanner.Text())
		}
	}()
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)
	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
			svc.logger.Info("Deploy Error: " + stderrScanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error : %v \n", err)
		svc.logger.Error("Collector Deploy Error", zap.Error(err))
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		svc.logger.Error("Collector Deploy Error", zap.Error(err))
	}

	return nil
}

func (svc maestroService) CreateOtelCollector(ctx context.Context, sinkID string, msg string, ownerID string) error {
	err := svc.collectorDeploy("apply", "otelcollectors", sinkID, k8sOtelCollector, otelCollectorCfg)
	if err != nil {
		return err
	}
	return nil
}

func (svc maestroService) UpdateOtelCollector(ctx context.Context, sinkID string, msg string, ownerID string) error {
	err := svc.collectorDeploy("apply", "otelcollectors", sinkID, k8sOtelCollector, otelCollectorCfg)
	if err != nil {
		return err
	}
	return nil
}

func (svc maestroService) DeleteOtelCollector(ctx context.Context, sinkID string, msg string, ownerID string) error {
	err := svc.collectorDeploy("delete", "otelcollectors", sinkID, k8sOtelCollector, otelCollectorCfg)
	if err != nil {
		return err
	}
	return nil
}
