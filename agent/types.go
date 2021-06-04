/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

type Sinks map[string]interface{}

type TLS struct {
	Verify bool `mapstructure:"verify"`
}

type OrbAgent struct {
	Backends []string          `mapstructure:"backends"`
	Vitals   map[string]string `mapstructure:"vitals"`
	MQTT     map[string]string `mapstructure:"mqtt"`
	TLS      TLS               `mapstructure:"tls"`
}

type Config struct {
	Debug    bool
	Version  float64  `mapstructure:"version"`
	Sinks    Sinks    `mapstructure:"sinks"`
	OrbAgent OrbAgent `mapstructure:"orb"`
}
