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
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/zap"
)

// Components returns of all imports and registering the `opentelemetry-collector-contrib` elements
func Components(logger zap.Logger) (component.Factories, error) {
	var err error
	factories := component.Factories{}
	extensions := getExtensions()
	factories.Extensions, err = component.MakeExtensionFactoryMap(extensions...)
	if err != nil {
		logger.Error("extension factories failure to load", zap.Error(err))
		return component.Factories{}, err
	}

	receivers := getReceivers()
	receivers = append(receivers, extraReceivers()...)
	factories.Receivers, err = component.MakeReceiverFactoryMap(receivers...)
	if err != nil {
		logger.Error("receivers factories failure to load", zap.Error(err))
		return component.Factories{}, err
	}

	exporters := getExporters()
	factories.Exporters, err = component.MakeExporterFactoryMap(exporters...)
	if err != nil {
		logger.Error("exporters factories failure to load", zap.Error(err))
		return component.Factories{}, err
	}

	processors := getProcessors()
	factories.Processors, err = component.MakeProcessorFactoryMap(processors...)
	if err != nil {
		logger.Error("processors factories failure to load", zap.Error(err))
		return component.Factories{}, err
	}

	return factories, nil
}

// getProcessors return processors factory, check version before adding and updating
func getProcessors() []component.ProcessorFactory {
	return []component.ProcessorFactory{
		// Inserts Tenant and Sinks data in otlp package
		// current version and stability for metrics [ 0.56.0 , alpha ]
		transformprocessor.NewFactory(),
	}
}

func getExporters() []component.ExporterFactory {
	return []component.ExporterFactory{
		// export log
		// current version and stability for metrics [ 0.56.0 , stable ]
		loggingexporter.NewFactory(),
		// current version and stability for metrics [ 0.56.0 , stable ]
		otlpexporter.NewFactory(),
		// current version and stability for metrics [ 0.56.0 , stable ]
		otlphttpexporter.NewFactory(),
		// current version and stability for metrics [ 0.56.0 , beta ]
		prometheusexporter.NewFactory(),
		// current version and stability for metrics [ 0.56.0 , beta ]
		prometheusremotewriteexporter.NewFactory(),
	}
}

func getReceivers() []component.ReceiverFactory {
	return []component.ReceiverFactory{
		// current version and stability for metrics [ 0.56.0 , stable ]
		otlpreceiver.NewFactory(),
	}
}

func getExtensions() []component.ExtensionFactory {
	return []component.ExtensionFactory{
		// not sure if we need that yet, this creates a ballast of memory
		// the ballast increases the base size of the heap so that our GC triggers are delayed and the number of GC
		//cycles over time is reduced
		// current version and stability [ 0.56.0 , beta ]
		ballastextension.NewFactory(),

		// Enables an extension that serves zPages, an HTTP endpoint that provides live data for debugging different
		// components that were properly instrumented for such.
		// All core exporters and receivers provide some zPage instrumentation.
		// current version and stability [ 0.56.0 , beta  ]
		zpagesextension.NewFactory(),
	}
}
