// Copyright 2019 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package components

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
)

func Components() (component.Factories, error) {
	var err error
	factories := component.Factories{}
	extensions := getExtensions()
	factories.Extensions, err = component.MakeExtensionFactoryMap(extensions...)
	if err != nil {
		return component.Factories{}, err
	}

	receivers := getReceivers()
	receivers = append(receivers, extraReceivers()...)
	factories.Receivers, err = component.MakeReceiverFactoryMap(receivers...)
	if err != nil {
		return component.Factories{}, err
	}

	exporters := getExporters()
	factories.Exporters, err = component.MakeExporterFactoryMap(exporters...)
	if err != nil {
		return component.Factories{}, err
	}

	processors := []component.ProcessorFactory{
		batchprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
	}
	factories.Processors, err = component.MakeProcessorFactoryMap(processors...)
	if err != nil {
		return component.Factories{}, err
	}

	return factories, nil
}

func getExporters() []component.ExporterFactory {
	return []component.ExporterFactory{
		loggingexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
	}
}

func getReceivers() []component.ReceiverFactory {
	return []component.ReceiverFactory{
		otlpreceiver.NewFactory(),
		prometheusreceiver.NewFactory(),
	}
}

func getExtensions() []component.ExtensionFactory {
	return []component.ExtensionFactory{
		ballastextension.NewFactory(),
		zpagesextension.NewFactory(),
	}
}
