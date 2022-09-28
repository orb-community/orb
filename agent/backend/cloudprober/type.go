package cloudprober

import "strings"

type Probes struct {
	probe []Probe `yaml:"probes"`
}

type Probe struct {
	name string `yaml:"p_name"`

	probeType string `yaml:"type"`

	targets string `yaml:"targets_host_names"`

	intervalMsec uint `yaml:"interval_msec"`

	timeoutMsec uint `yaml:"timeout_msec"`
}

func (p *Probes) ToConfigFile() string {
	var builder strings.Builder
	for _, probeCfg := range p.probe {
		builder.WriteString("\nprobe { \n")
		builder.WriteString("  name: \"" + probeCfg.name + "\"\n")
		builder.WriteString("  type: " + probeCfg.probeType + "\n")
		builder.WriteString("  targets { \n")
		builder.WriteString("    host_names: ")
		builder.WriteString(probeCfg.targets)
		builder.WriteString("\n")
		builder.WriteString("  interval_msec: ")
		builder.WriteRune(int32(probeCfg.intervalMsec))
		builder.WriteString("\n")
		builder.WriteString("  timeout_msec: ")
		builder.WriteRune(int32(probeCfg.timeoutMsec))
		builder.WriteString("\n")
		builder.WriteString("}")
		builder.WriteString("\n")
	}
	return builder.String()
}
