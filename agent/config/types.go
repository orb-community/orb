/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

type TLS struct {
	Verify bool `mapstructure:"verify"`
}

type OrbAPIConfig struct {
	Address string
	Token   string
}

type MQTTConfig struct {
	Address   string
	Id        string
	Key       string
	ChannelID string
}

type OrbAgent struct {
	Backends map[string]map[string]string `mapstructure:"backends"`
	Tags     map[string]string            `mapstructure:"tags"`
	Cloud    map[string]map[string]string `mapstructure:"cloud"`
	TLS      TLS                          `mapstructure:"tls"`
	DB       map[string]string            `mapstructure:"db"`
}

type Config struct {
	Debug    bool
	Version  float64  `mapstructure:"version"`
	OrbAgent OrbAgent `mapstructure:"orb"`
}
