/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package promremotewrite

type PromRemoteConfig struct {
}

type promRemoteWrite struct {
}

func (p promRemoteWrite) Close() error {
	panic("implement me")
}

type PromRemoteWriter interface {
	Close() error
}

// New instantiates the prom remote write
func New(config PromRemoteConfig) PromRemoteWriter {
	return &promRemoteWrite{}
}
