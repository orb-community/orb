/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package prometheus

import (
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks/backend"
)

var _ backend.Backend = (*prometheusBackend)(nil)

type prometheusBackend struct {
}

func (p prometheusBackend) Validate(config types.Metadata) error {
	return nil
}

func Register() bool {
	backend.Register("prometheus", &prometheusBackend{})
	return true
}
