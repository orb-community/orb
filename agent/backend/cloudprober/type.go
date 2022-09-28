package cloudprober

import (
	"strconv"
	"strings"
)

/**
{
    "name": "policy_hackathon",
    "description": "some_description",
    "backend": "cloudprober",
    "tags": {
        "key": "value"
    },
    "policy": {
        "probes": [
            {
				"name" : "p_name"
				"type": "http",
				"interval_msec": 2000,
				"timeout_msec": 3000,
				"targets_host_names": "www.google.com,ns1.com"
			}
        ]
    }
}
*/

type Probes struct {
	ProbeData []ProbeData `json:"probes" mapstructure:"probe"`
}

type ProbeData struct {
	Name         string `json:"name" mapstructure:"name"`
	ProbeType    string `json:"type" mapstructure:"type"`
	Targets      string `json:"targets_host_names" mapstructure:"targets_host_names"`
	IntervalMsec uint   `json:"interval_msec" mapstructure:"interval_msec"`
	TimeoutMsec  uint   `json:"timeout_msec" mapstructure:"timeout_msec"`
}

func (p *Probes) ToConfigFile() string {
	var builder strings.Builder
	for _, probeCfg := range p.ProbeData {
		builder.WriteString("\nprobe { \n")
		builder.WriteString("  name: \"" + probeCfg.Name + "\"\n")
		builder.WriteString("  type: " + probeCfg.ProbeType + "\n")
		builder.WriteString("  targets { \n")
		builder.WriteString("    host_names: ")
		builder.WriteString(probeCfg.Targets)
		builder.WriteString("\n")
		builder.WriteString("  interval_msec: ")
		builder.WriteString(strconv.Itoa(int(probeCfg.IntervalMsec)))
		builder.WriteString("\n")
		builder.WriteString("  timeout_msec: ")
		builder.WriteString(strconv.Itoa(int(probeCfg.TimeoutMsec)))
		builder.WriteString("\n")
		builder.WriteString("}")
		builder.WriteString("\n")
	}
	return builder.String()
}
