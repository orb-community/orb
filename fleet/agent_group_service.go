// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/pkg/errors"
)

var (
	ErrCreateAgentGroup = errors.New("failed to create agent group")
)

func (svc fleetService) CreateAgentGroup(ctx context.Context, token string, s AgentGroup) (AgentGroup, error) {
	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return AgentGroup{}, err
	}

	s.MFOwnerID = mfOwnerID

	md := map[string]interface{}{"type": "orb_agent_group"}

	// create main Agent RPC Channel
	mfChannelID, err := svc.mfsdk.CreateChannel(mfsdk.Channel{
		Name:     s.Name.String(),
		Metadata: md,
	}, token)
	if err != nil {
		return AgentGroup{}, errors.Wrap(ErrCreateAgent, err)
	}

	s.MFChannelID = mfChannelID

	err = svc.agentGroupRepository.Save(ctx, s)
	if err != nil {
		return AgentGroup{}, errors.Wrap(ErrCreateAgentGroup, err)
	}

	return s, nil
}
