// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/ns1labs/orb/fleet/pb"
	"go.uber.org/zap"
)

const publisher = "orb_policies"

type PolicyCommsService interface {
}

var _ PolicyCommsService = (*policiesCommsService)(nil)

type policiesCommsService struct {
	logger      *zap.Logger
	fleetClient pb.FleetServiceClient
}

func NewPoliciesCommsService(logger *zap.Logger, fleetClient pb.FleetServiceClient, policyPubSub mfnats.PubSub) PolicyCommsService {
	return &policiesCommsService{
		logger:      logger,
		fleetClient: fleetClient,
	}
}
