/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

type Pktvisor map[string]interface{}

type Sinks map[string]interface{}

type OrbAgent struct {
	Vitals map[string]string `mapstructure:"vitals"`
	MQTT   map[string]string `mapstructure:"mqtt"`
}

type Config struct {
	Version  float64  `mapstructure:"version"`
	Pktvisor Pktvisor `mapstructure:"pktvisor"`
	Sinks    Sinks    `mapstructure:"sinks"`
	OrbAgent OrbAgent `mapstructure:"orb"`
}
