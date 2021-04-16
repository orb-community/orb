package agent

import (
	"fmt"
	"github.com/ns1labs/orb/datasource"
)

type Pktvisor map[string]interface{}

type Sinks map[string]interface{}

type OrbAgent struct {
	Vitals     map[string]string       `mapstructure:"vitals"`
	Datasource datasource.ConsulConfig `mapstructure:"datasource"`
}

type Config struct {
	Version  float64  `mapstructure:"version"`
	Pktvisor Pktvisor `mapstructure:"pktvisor"`
	Sinks    Sinks    `mapstructure:"sinks"`
	OrbAgent OrbAgent `mapstructure:"orb"`
}

type Agent struct {
}

func New(c Config) (*Agent, error) {
	return &Agent{}, nil
}

func (a *Agent) Start() error {
	fmt.Println("started")
	return nil
}
