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
	ID          string
	Name        types.Identifier
	Description string
	MFOwnerID   string
	Backend     string
	Version     int32
	OrbTags     types.Tags
	Policy      types.Metadata
	Created     time.Time
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

type Page struct {
	PageMetadata
	Policies []Policy
}

type Service interface {
	// AddPolicy creates new agent Policy
	AddPolicy(ctx context.Context, token string, p Policy, format string, policyData string) (Policy, error)

	// ViewPolicyByID retrieving policy by id with token
	ViewPolicyByID(ctx context.Context, token string, policyID string) (Policy, error)

	// ListPolicies
	ListPolicies(ctx context.Context, token string, pm PageMetadata) (Page, error)

	// ViewPolicyByIDInternal gRPC version of retrieving policy by id with no token
	ViewPolicyByIDInternal(ctx context.Context, policyID string, ownerID string) (Policy, error)

	// ListPoliciesByGroupIDInternal gRPC version of retrieving list of policies belonging to specified agent group with no token
	ListPoliciesByGroupIDInternal(ctx context.Context, groupIDs []string, ownerID string) ([]Policy, error)

	// EditPolicy edit a existing policy by id with a valid token
	EditPolicy(ctx context.Context, token string, pol Policy, format string, policyData string) (Policy, error)

	// AddDataset creates new Dataset
	AddDataset(ctx context.Context, token string, d Dataset) (Dataset, error)

	// InactivateDatasetByGroupID
	InactivateDatasetByGroupID(ctx context.Context, groupID string, token string) error

	// ListDatasetsByPolicyID retrieves the subset of Datasets by policyID owned by the specified user
	ListDatasetsByPolicyIDInternal(ctx context.Context, policyID string, token string) ([]Dataset, error)
}

type Repository interface {
	// SavePolicy persists a Policy. Successful operation is indicated by non-nil
	// error response.
	SavePolicy(ctx context.Context, policy Policy) (string, error)

	// RetrievePolicyByID Retrieve policy by id
	RetrievePolicyByID(ctx context.Context, policyID string, ownerID string) (Policy, error)

	// RetrievePoliciesByGroupID Retrieve policy list by group id
	RetrievePoliciesByGroupID(ctx context.Context, groupIDs []string, ownerID string) ([]Policy, error)

	// RetrieveAll retrieves the subset of Policies owned by the specified user
	RetrieveAll(ctx context.Context, owner string, pm PageMetadata) (Page, error)

	// UpdatePolicy update a existing policy by id with a valid token
	UpdatePolicy(ctx context.Context, owner string, pol Policy) error

	// SaveDataset persists a Dataset. Successful operation is indicated by non-nil
	// error response.
	SaveDataset(ctx context.Context, dataset Dataset) (string, error)

	// InactivateDatasetByGroupID inactivate a dataset
	InactivateDatasetByGroupID(ctx context.Context, groupID string, ownerID string) error

	// RetrieveDatasetsByPolicyID retrieves the subset of Datasets by policyID owned by the specified user
	RetrieveDatasetsByPolicyID(ctx context.Context, policyID string, ownerID string) ([]Dataset, error)
}
