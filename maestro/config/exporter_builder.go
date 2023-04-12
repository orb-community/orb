package config

import "github.com/orb-community/orb/pkg/types"

type ExporterConfigService interface {
	GetExportersFromMetadata(config types.Metadata, authenticationExtensionName string) (Exporters, string)
}

func FromStrategy(backend string) ExporterConfigService {
	switch backend {
	case "prometheus":
		return &PrometheusExporterConfig{}
	case "otlpexporter":
		return &OTLPExporterBuilder{}
	}

	return nil
}

type PrometheusExporterConfig struct {
}

func (p *PrometheusExporterConfig) GetExportersFromMetadata(config types.Metadata, authenticationExtensionName string) (Exporters, string) {
	endpointCfg := config.GetSubMetadata("exporter")["remote_host"].(string)
	return Exporters{
		PrometheusRemoteWrite: &PrometheusRemoteWriteExporterConfig{
			Endpoint: endpointCfg,
			Auth:     Auth{Authenticator: authenticationExtensionName},
		},
	}, "prometheusremotewrite"
}

type OTLPExporterBuilder struct {
}

func (O *OTLPExporterBuilder) GetExportersFromMetadata(config types.Metadata, authenticationExtensionName string) (Exporters, string) {
	endpointCfg := config.GetSubMetadata("exporter")["endpoint"].(string)
	return Exporters{
		OTLPExporter: &OTLPExporterConfig{
			Endpoint: endpointCfg,
			Auth:     Auth{Authenticator: authenticationExtensionName},
		},
	}, "otlphttp"
}
