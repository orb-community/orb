/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"errors"
	"github.com/ns1labs/orb/agent/backend"
	"go.uber.org/zap"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	logger *zap.Logger
	binary string
}

func (p *pktvisorBackend) Start() error {
	p.logger.Info("pktvisor starting")
	return nil
}

func (p *pktvisorBackend) Stop() error {
	p.logger.Info("pktvisor stopping")
	return nil
}

func (p *pktvisorBackend) Configure(logger *zap.Logger, config map[string]string) error {
	p.logger = logger

	var prs bool
	if p.binary, prs = config["binary"]; !prs {
		return errors.New("you must specify pktvisor binary")
	}
	return nil
}

func Register() bool {
	backend.Register("pktvisor", &pktvisorBackend{})
	return true
}
