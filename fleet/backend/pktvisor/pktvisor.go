/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"encoding/json"
	"fmt"
	"github.com/ns1labs/orb/fleet/backend"
	"github.com/ns1labs/orb/pkg/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	Backend     string
	Description string
}

func (p pktvisorBackend) Metadata() interface{} {
	return p
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{
		Backend:     "pktvisor",
		Description: "pktvisor observability agent from pktvisor.dev",
	})
	return true
}

func (p pktvisorBackend) Handlers() (_ types.Metadata, err error) {
	wd := getWorkDirectory()
	jsonFile, err := ioutil.ReadFile(fmt.Sprintf("%s/fleet/backend/pktvisor/handlers.json", wd))
	if err != nil {
		return nil, err
	}
	var handlers types.Metadata
	err = json.Unmarshal([]byte(jsonFile), &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}

func (p pktvisorBackend) Inputs() (_ types.Metadata, err error) {
	wd := getWorkDirectory()
	jsonFile, err := ioutil.ReadFile(fmt.Sprintf("%s/fleet/backend/pktvisor/inputs.json", wd))
	if err != nil {
		return nil, err
	}
	var handlers types.Metadata
	err = json.Unmarshal([]byte(jsonFile), &handlers)
	if err != nil {
		return nil, err
	}
	return handlers, nil
}

func getWorkDirectory() string {
	// When you works with tests, the path it's different from the prod running
	// So here I'm getting the right working directory, no matter if its test or prod
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "orb") {
		wd = filepath.Dir(wd)
	}
	return wd
}
