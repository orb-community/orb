/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"encoding/json"
	"github.com/ns1labs/orb/fleet/backend"
	"github.com/ns1labs/orb/pkg/types"
	"io/ioutil"
	"os"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	Backend     string
	Description string
}

func (p pktvisorBackend) Metadata() interface{} {
	return pktvisorBackend{
		Backend:     "pktvisor",
		Description: "pktvisor observability agent from pktvisor.dev",
	}
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{})
	return true
}

func (p pktvisorBackend) Handlers() (_ types.Metadata, err error) {
	jsonFile, err := os.Open("fleet/backend/pktvisor/handlers.json")
	if err != nil {
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var handlers types.Metadata
	err = json.Unmarshal(byteValue, &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}

func (p pktvisorBackend) Inputs() (_ types.Metadata, err error) {
	jsonFile, err := os.Open("fleet/backend/pktvisor/inputs.json")
	if err != nil {
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var handlers types.Metadata
	err = json.Unmarshal(byteValue, &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}
