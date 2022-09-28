package cloudprober

import "strings"

type Probes struct {
	probe []Probe
}

type Probe struct {
	name string `yaml:"name"`

	probeType string `yaml:"type"`

	targets target `yaml:"targets"`

	intervalMsec uint `yaml:"interval_msec"`

	timeoutMsec uint `yaml:"timeout_msec"`
}

type target struct {
	hostNames []string `yaml:"host_names"`
}

func (p *Probes) ToConfigFile() string {
	var builder strings.Builder
	for _, probeCfg := range p.probe {
		builder.WriteString("\nprobe { \n")
		builder.WriteString("  name: \"" + probeCfg.name + "\"\n")
		builder.WriteString("  type: " + probeCfg.probeType + "\n")
		builder.WriteString("  targets { \n")
		builder.WriteString("    host_names: ")
		builder.WriteString(strings.Join(probeCfg.targets.hostNames, ","))
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
