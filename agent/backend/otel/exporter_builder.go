package otel

// TODO Create a struct to hold the exporter and processor as default to inject the policy id and name as attribute with an attribute processor
type defaultExporters struct {
	Otlp *defaultOtlpExporter `yaml:"otlp"`
}

type defaultOtlpExporter struct {
}

func (o *openTelemetryBackend) buildDefaultExporterAndProcessor() {

}
