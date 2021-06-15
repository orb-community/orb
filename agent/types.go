/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

type OrbAgentInfo struct {
	Version string `json:"version"`
}

type BackendInfo struct {
	Version string `json:"version"`
}

type TLS struct {
	Verify bool `mapstructure:"verify"`
}

type OrbAgent struct {
	Backends map[string]map[string]string `mapstructure:"backends"`
	Tags     map[string]string            `mapstructure:"tags"`
	MQTT     map[string]string            `mapstructure:"mqtt"`
	TLS      TLS                          `mapstructure:"tls"`
}

type Config struct {
	Debug    bool
	Version  float64  `mapstructure:"version"`
	OrbAgent OrbAgent `mapstructure:"orb"`
}
