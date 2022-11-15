/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package config

type ConfigRepo interface {
	Exists(ownerID string, sinkID string) bool
	Add(config SinkConfig) error
	Remove(ownerID string, sinkID string) error
	Get(ownerID string, sinkID string) (SinkConfig, error)
	Edit(config SinkConfig) error
	GetAll(ownerID string) ([]SinkConfig, error)
	GetAllOwners() ([]string, error)
}