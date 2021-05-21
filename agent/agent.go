/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"fmt"
)

type Pktvisor map[string]interface{}

type Sinks map[string]interface{}

type OrbAgent struct {
	Vitals map[string]string `mapstructure:"vitals"`
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
