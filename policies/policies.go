/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package policies

import (
	"context"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

type Policy struct {
	ID        string
	Name      types.Identifier
	MFOwnerID string
	Backend   string
	Version   int32
	OrbTags   types.Tags
	Policy    types.Metadata
	Created   time.Time
}

type Dataset struct {
	ID           string
	Name         types.Identifier
	MFOwnerID    string
	Valid        bool
	AgentGroupID string
	PolicyID     string
	SinkID       string
	Metadata     types.Metadata
	Created      time.Time
}

type Repository interface {
	// SavePolicy persists a Policy. Successful operation is indicated by non-nil
	// error response.
	SavePolicy(ctx context.Context, policy Policy) (string, error)

	// RetrievePolicyByID Retrieve policy by id
	RetrievePolicyByID(ctx context.Context, policyID string, ownerID string) (Policy, error)

	// RetrievePoliciesByGroupID Retrieve policy list by group id
	RetrievePoliciesByGroupID(ctx context.Context, groupIDs []string, ownerID string) ([]Policy, error)

	// SaveDataset persists a Dataset. Successful operation is indicated by non-nil
	// error response.
	SaveDataset(ctx context.Context, dataset Dataset) (string, error)

	// InactivateDatasetByGroupID inactivate a dataset
	InactivateDatasetByGroupID(ctx context.Context, groupID string, ownerID string) error
}
