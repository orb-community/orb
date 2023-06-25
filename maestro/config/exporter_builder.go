package config

import (
	"log"

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
    exporters := Exporters{}
    exporterMetadata := config.GetSubMetadata("exporter")
    if exporterMetadata == nil {
        log.Println("exporter metadata is missing")
        return exporters, ""
    }
    endpointCfg, ok := exporterMetadata["remote_host"].(string)
    if !ok {
		log.Printf("remote_host metadata is missing or not a string %v", exporterMetadata)
        return exporters, ""
    }
	exporters.PrometheusRemoteWrite = &PrometheusRemoteWriteExporterConfig{
		Endpoint: endpointCfg,
		Auth:     Auth{Authenticator: authenticationExtensionName},
	}
	// Check to add X-Scope-OrgID header
	header := config.GetSubMetadata("headers")
	if header != nil {
		headerStr := header["X-Scope-OrgID"].(string)
		if headerStr != "" {
			exporters.OTLPExporter.Headers = map[string]string{
				"X-Scope-OrgID": headerStr,
			}
		}
	}

    return exporters, "prometheusremotewrite"
}




type OTLPHTTPExporterBuilder struct {
}

func (O *OTLPHTTPExporterBuilder) GetExportersFromMetadata(config types.Metadata, authenticationExtensionName string) (Exporters, string) {
    exporters := Exporters{}
    exporterMetadata := config.GetSubMetadata("exporter")
    if exporterMetadata == nil {
        log.Println("exporter metadata is missing")
        return exporters, ""
    }
    endpointCfg, ok := exporterMetadata["endpoint"].(string)
    if !ok {
		log.Printf("endpoint metadata is missing or not a string %v", exporterMetadata)
        return exporters, ""
    }
    exporters.OTLPExporter = &OTLPExporterConfig{
        Endpoint: endpointCfg,
        Auth:     Auth{Authenticator: authenticationExtensionName},
    }
	// Check to add X-Scope-OrgID header
	header := config.GetSubMetadata("headers")
	if header != nil {
		headerStr := header["X-Scope-OrgID"].(string)
		if headerStr != "" {
			exporters.OTLPExporter.Headers = map[string]string{
				"X-Scope-OrgID": headerStr,
			}
		}
	}
    return exporters, "otlphttp"
}
