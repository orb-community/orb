package config

import "github.com/orb-community/orb/pkg/types"

type ExporterConfigService interface {
	GetExportersFromMetadata(config types.Metadata, authenticationExtensionName string) (Exporters, string)
}

func FromStrategy(backend string) ExporterConfigService {
	switch backend {
	case "prometheus":
		return &PrometheusExporterConfig{}
	case "otlphttp":
		return &OTLPHTTPExporterBuilder{}
	}

	return nil
}

type PrometheusExporterConfig struct {
}

func (p *PrometheusExporterConfig) GetExportersFromMetadata(config types.Metadata, authenticationExtensionName string) (Exporters, string) {
	exporterSubMeta := config.GetSubMetadata("exporter")
	endpointCfg := exporterSubMeta["remote_host"].(string)
	customHeaders := exporterSubMeta["headers"]
	if customHeaders == nil {
		return Exporters{
			PrometheusRemoteWrite: &PrometheusRemoteWriteExporterConfig{
				Endpoint: endpointCfg,
				Auth:     Auth{Authenticator: authenticationExtensionName},
			},
		}, "prometheusremotewrite"
	}
	return Exporters{
		PrometheusRemoteWrite: &PrometheusRemoteWriteExporterConfig{
			Endpoint: endpointCfg,
			Auth:     Auth{Authenticator: authenticationExtensionName},
			Headers:  customHeaders.(map[string]interface{}),
		},
	}, "prometheusremotewrite"
}

type OTLPHTTPExporterBuilder struct {
}

func (O *OTLPHTTPExporterBuilder) GetExportersFromMetadata(config types.Metadata, authenticationExtensionName string) (Exporters, string) {
	endpointCfg := config.GetSubMetadata("exporter")["endpoint"].(string)
	return Exporters{
		OTLPExporter: &OTLPExporterConfig{
			Endpoint: endpointCfg,
			Auth:     Auth{Authenticator: authenticationExtensionName},
		},
	}, "otlphttp"
}
