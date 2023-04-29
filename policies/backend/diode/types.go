/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package diode

import "github.com/orb-community/orb/pkg/types"

const CurrentSchemaVersion = "1.0"

// policy for suzieq
type collectionPolicy struct {
	Backend string         `json:"backend"`
	Data    types.Metadata `json:"data"`
	Kind    string         `json:"kind"`
}

/*

kind: discovery
backend: suzieq
config:
  netbox:
    site: New York NY
data:
  inventory:
    sources:
      - name: default_inventory
  hosts:
    - url: ssh://1.2.3.4:2021 username=user password=password
    - url: ssh://resolvable.host.name username=user password=password
  devices:
    - name: default_devices
      transport: ssh
      ignore-known-hosts: true
      slow-host: true
  namespaces:
    - name: default_namespace
      source: default_inventory
      device: default_devices

*/
