package output

import (
	"github.com/orb-community/orb/pkg/types"
)

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
	if exporterSubMeta == nil {
		return Exporters{}, ""
	}
	endpointCfg, ok := exporterSubMeta["remote_host"].(string)
	if !ok {
		return Exporters{}, ""
	}
	var auth *Auth
	if authenticationExtensionName != "" {
		auth = &Auth{Authenticator: authenticationExtensionName}
	}

	customHeaders, ok := exporterSubMeta["headers"]
	if !ok || customHeaders == nil {
		return Exporters{
			PrometheusRemoteWrite: &PrometheusRemoteWriteExporterConfig{
				Endpoint: endpointCfg,
				Auth:     auth,
			},
		}, "prometheusremotewrite"
	}
	return Exporters{
		PrometheusRemoteWrite: &PrometheusRemoteWriteExporterConfig{
			Endpoint: endpointCfg,
			Auth:     auth,
			Headers:  customHeaders.(map[string]interface{}),
		},
	}, "prometheusremotewrite"
}

type OTLPHTTPExporterBuilder struct {
}

func (O *OTLPHTTPExporterBuilder) GetExportersFromMetadata(config types.Metadata, authenticationExtensionName string) (Exporters, string) {
	exporterSubMeta := config.GetSubMetadata("exporter")
	endpointCfg := exporterSubMeta["endpoint"].(string)
	var auth *Auth
	if authenticationExtensionName != "" {
		auth = &Auth{Authenticator: authenticationExtensionName}
	}
	customHeaders, ok := exporterSubMeta["headers"]
	if !ok || customHeaders == nil {
		return Exporters{
			OTLPExporter: &OTLPExporterConfig{
				Endpoint: endpointCfg,
				Auth:     auth,
			},
		}, "otlphttp"
	} else {
		return Exporters{
			OTLPExporter: &OTLPExporterConfig{
				Endpoint: endpointCfg,
				Auth:     auth,
				Headers:  customHeaders.(map[string]interface{}),
			},
		}, "otlphttp"
	}
}
