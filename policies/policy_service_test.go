package policies_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mainflux/mainflux"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	policies "github.com/ns1labs/orb/policies"
	plmocks "github.com/ns1labs/orb/policies/mocks"
	sinkmocks "github.com/ns1labs/orb/sinks/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	token        = "token"
	invalidToken = "invalid"
	email        = "user@example.com"
	format       = "yaml"
	policy_data  = `handlers:
  modules:
    default_dns:
      type: dns
    default_net:
      type: net
input:
  input_type: pcap
  tap: default_pcap
kind: collection`
	limit   = 10
	wrongID = "28ea82e7-0224-4798-a848-899a75cdc650"
)

func newService(auth mainflux.AuthServiceClient) policies.Service {
	policyRepo := plmocks.NewPoliciesRepository()
	fleetGrpcClient := flmocks.NewClient()
	SinkServiceClient := sinkmocks.NewClient()

	return policies.New(nil, auth, policyRepo, fleetGrpcClient, SinkServiceClient)
}

func TestRetrievePolicyByID(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	policy := createPolicy(t, svc, "policy")

	cases := map[string]struct {
		id    string
		token string
		err   error
	}{
		"view a existing policy": {
			id:    policy.ID,
			token: token,
			err:   nil,
		},
		"view policy with wrong credentials": {
			id:    policy.ID,
			token: "wrong",
			err:   policies.ErrUnauthorizedAccess,
		},
		"view non-existing policy": {
			id:    "9bb1b244-a199-93c2-aa03-28067b431e2c",
			token: token,
			err:   policies.ErrNotFound,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.ViewPolicyByID(context.Background(), tc.token, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestListPolicies(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	var policyList []policies.Policy
	for i := 0; i < limit; i++ {
		pl := createPolicy(t, svc, fmt.Sprintf("policy-%d", i))
		policyList = append(policyList, pl)
	}

	cases := map[string]struct {
		token string
		pm    policies.PageMetadata
		size  uint64
		err   error
	}{
		"retrieve a list of policies": {
			token: token,
			pm: policies.PageMetadata{
				Limit:  limit,
				Offset: 0,
			},
			size: limit,
			err:  nil,
		},
		"list half": {
			token: token,
			pm: policies.PageMetadata{
				Offset: limit / 2,
				Limit:  limit,
			},
			size: limit / 2,
			err:  nil,
		},
		"list last policy": {
			token: token,
			pm: policies.PageMetadata{
				Offset: limit - 1,
				Limit:  limit,
			},
			size: 1,
			err:  nil,
		},
		"list empty set": {
			token: token,
			pm: policies.PageMetadata{
				Offset: limit + 1,
				Limit:  limit,
			},
			size: 0,
			err:  nil,
		},
		"list with zero limit": {
			token: token,
			pm: policies.PageMetadata{
				Offset: 1,
				Limit:  0,
			},
			size: 0,
			err:  nil,
		},
		"list with wrong credentials": {
			token: "wrong",
			pm: policies.PageMetadata{
				Offset: 0,
				Limit:  0,
			},
			size: 0,
			err:  policies.ErrUnauthorizedAccess,
		},
		"list all policies sorted by name ascendent": {
			token: token,
			pm: policies.PageMetadata{
				Offset: 0,
				Limit:  limit,
				Order:  "name",
				Dir:    "asc",
			},
			size: limit,
			err:  nil,
		},
		"list all policies sorted by name descendent": {
			token: token,
			pm: policies.PageMetadata{
				Offset: 0,
				Limit:  limit,
				Order:  "name",
				Dir:    "desc",
			},
			size: limit,
			err:  nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			page, err := svc.ListPolicies(context.Background(), tc.token, tc.pm)
			size := uint64(len(page.Policies))
			assert.Equal(t, size, tc.size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, size))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
			testSortPolicies(t, tc.pm, page.Policies)
		})

	}
}

func TestEditPolicy(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	policy := createPolicy(t, svc, "policy")

	nameID, err := types.NewIdentifier("new-policy")
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	wrongOwnerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	wrongPolicy := policies.Policy{MFOwnerID: wrongOwnerID.String()}
	newPolicy := policies.Policy{
		ID:         policy.ID,
		Name:       nameID,
		MFOwnerID:  policy.MFOwnerID,
		PolicyData: policy_data,
		Format:     format,
	}

	invalidFormatPolicy := newPolicy
	invalidFormatPolicy.Format = "invalid"

	invalidPolicyData := newPolicy
	invalidPolicyData.PolicyData = "invalid"

	cases := map[string]struct {
		pol   policies.Policy
		token string
		err   error
	}{
		"update a existing policy": {
			pol:   newPolicy,
			token: token,
			err:   nil,
		},
		"update policy with wrong credentials": {
			pol:   newPolicy,
			token: "invalidToken",
			err:   policies.ErrUnauthorizedAccess,
		},
		"update a non-existing policy": {
			pol:   wrongPolicy,
			token: token,
			err:   policies.ErrNotFound,
		},
		"update a existing policy with invalid format": {
			pol:   invalidFormatPolicy,
			token: token,
			err:   policies.ErrValidatePolicy,
		},
		"update a existing policy with invalid policy_data": {
			pol:   invalidPolicyData,
			token: token,
			err:   policies.ErrValidatePolicy,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			res, err := svc.EditPolicy(context.Background(), tc.token, tc.pol)
			if err == nil {
				assert.Equal(t, tc.pol.Name.String(), res.Name.String(), fmt.Sprintf("%s: expected name %s got %s", desc, tc.pol.Name.String(), res.Name.String()))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected error %d got %d", desc, tc.err, err))
		})
	}

}

func TestRemovePolicy(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	plcy := createPolicy(t, svc, "policy")

	cases := map[string]struct {
		id    string
		token string
		err   error
	}{
		"Remove a existing policy": {
			id:    plcy.ID,
			token: token,
			err:   nil,
		},
		"delete non-existent policy": {
			id:    wrongID,
			token: token,
			err:   nil,
		},
		"delete policy with wrong credentials": {
			id:    plcy.ID,
			token: invalidToken,
			err:   policies.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := svc.RemovePolicy(context.Background(), tc.token, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestValidatePolicy(t *testing.T) {
	var nameID, _ = types.NewIdentifier("my-policy")
	var policy = policies.Policy{
		Name:       nameID,
		Backend:    "pktvisor",
		OrbTags:    map[string]string{"region": "eu", "node_type": "dns"},
		Format:     format,
		PolicyData: policy_data,
	}

	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	cases := map[string]struct {
		policy policies.Policy
		token  string
		err    error
	}{
		"validate a new policy": {
			policy: policy,
			token:  token,
			err:    nil,
		},
		"validate a policy with a invalid token": {
			policy: policy,
			token:  invalidToken,
			err:    policies.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.ValidatePolicy(context.Background(), tc.token, policy)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}

}

func TestCreatePolicy(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpect error: %s", err))

	nameID, _ := types.NewIdentifier("my-policy")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policy := policies.Policy{
		Name:        nameID,
		MFOwnerID:   ownerID.String(),
		Description: "An example policy",
		Backend:     "pktvisor",
		Version:     0,
		OrbTags:     map[string]string{"region": "eu"},
		PolicyData:  policy_data,
		Format:      format,
	}

	cases := map[string]struct {
		policy policies.Policy
		token  string
		err    error
	}{
		"create a new policy": {
			policy: policy,
			token:  token,
			err:    nil,
		},
		"create a policy with an invalid token": {
			policy: policy,
			token:  invalidToken,
			err:    policies.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.AddPolicy(context.Background(), tc.token, tc.policy)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, err, tc.err))
			t.Log(tc.token)
		})
	}
}

func TestCreateDataset(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpect error: %s", err))

	nameID, _ := types.NewIdentifier("my-policy")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:      nameID,
		MFOwnerID: ownerID.String(),
	}

	cases := map[string]struct {
		dataset policies.Dataset
		token   string
		err     error
	}{
		"create a new dataset": {
			dataset: dataset,
			token:   token,
			err:     nil,
		},
		"create a dataset with an invalid token": {
			dataset: dataset,
			token:   invalidToken,
			err:     policies.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.AddDataset(context.Background(), tc.token, tc.dataset)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, err, tc.err))
			t.Log(tc.token)
		})
	}
}

func TestEditDataset(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	policy := createPolicy(t, svc, "policy")

	groupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	sinkID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	validName, err := types.NewIdentifier("dataset")
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset := policies.Dataset{
		Name:         validName,
		Valid:        true,
		AgentGroupID: groupID.String(),
		PolicyID:     policy.ID,
		SinkIDs:      []string{sinkID.String()},
	}

	dataset, err = svc.AddDataset(context.Background(), token, dataset)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	nameID, err := types.NewIdentifier("new-policy")
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	wrongOwnerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	wrongSinkID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	wrongSinkDs := createDataset(t, svc, "wrong_sink")
	wrongSinkDs.SinkIDs = []string{wrongSinkID.String()}

	wrongDataset := policies.Dataset{MFOwnerID: wrongOwnerID.String()}
	newDataset := policies.Dataset{
		ID:        dataset.ID,
		Name:      nameID,
		MFOwnerID: dataset.MFOwnerID,
		SinkIDs:   dataset.SinkIDs,
	}

	cases := map[string]struct {
		ds    policies.Dataset
		token string
		err   error
	}{
		"update a existing dataset": {
			ds:    newDataset,
			token: token,
			err:   nil,
		},
		"update dataset with wrong credentials": {
			ds:    newDataset,
			token: "invalidToken",
			err:   policies.ErrUnauthorizedAccess,
		},
		"update a non-existing dataset": {
			ds:    wrongDataset,
			token: token,
			err:   policies.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			res, err := svc.EditDataset(context.Background(), tc.token, tc.ds)
			if err == nil {
				assert.Equal(t, tc.ds.Name.String(), res.Name.String(), fmt.Sprintf("%s: expected name %s got %s", desc, tc.ds.Name.String(), res.Name.String()))
			}
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected error %d got %d", desc, tc.err, err))
		})
	}
}

func TestRemoveDataset(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	ds := createDataset(t, svc, "dataset")

	cases := map[string]struct {
		id    string
		token string
		err   error
	}{
		"Remove a existing dataset": {
			id:    ds.ID,
			token: token,
			err:   nil,
		},
		"delete non-existent dataset": {
			id:    wrongID,
			token: token,
			err:   nil,
		},
		"delete dataset with wrong credentials": {
			id:    ds.ID,
			token: invalidToken,
			err:   policies.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := svc.RemoveDataset(context.Background(), tc.token, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestValidateDataset(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	policy := createPolicy(t, svc, "policy")
	var nameID, _ = types.NewIdentifier("my-dataset")
	var (
		sinkIDsArray               = []string{"f5b2d342-211d-a9ab-1233-63199a3fc16f", "03679425-aa69-4574-bf62-e0fe71b80939"}
		dataset                    = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: policy.ID, SinkIDs: sinkIDsArray, Valid: true}
		datasetEmptySinkID         = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: policy.ID, SinkIDs: []string{}, Valid: true}
		datasetEmptyPolicyID       = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: "", SinkIDs: sinkIDsArray, Valid: true}
		datasetEmptyAgentGroupID   = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "", PolicyID: policy.ID, SinkIDs: sinkIDsArray, Valid: true}
		datasetInvalidSinkID       = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: policy.ID, SinkIDs: []string{"invalid"}, Valid: true}
		datasetInvalidPolicyID     = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: "invalid", SinkIDs: sinkIDsArray, Valid: true}
		datasetInvalidAgentGroupID = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "invalid", PolicyID: policy.ID, SinkIDs: sinkIDsArray, Valid: true}
	)

	cases := map[string]struct {
		dataset policies.Dataset
		token   string
		err     error
	}{
		"validate a new dataset": {
			dataset: dataset,
			token:   token,
			err:     nil,
		},
		"validate a dataset with a invalid token": {
			dataset: dataset,
			token:   invalidToken,
			err:     policies.ErrUnauthorizedAccess,
		},
		"validate a dataset with a empty sink ID": {
			dataset: datasetEmptySinkID,
			token:   token,
			err:     policies.ErrMalformedEntity,
		},
		"validate a dataset with a empty policy ID": {
			dataset: datasetEmptyPolicyID,
			token:   token,
			err:     policies.ErrMalformedEntity,
		},
		"validate a dataset with a empty agent group ID": {
			dataset: datasetEmptyAgentGroupID,
			token:   token,
			err:     policies.ErrMalformedEntity,
		},
		"validate a dataset with a invalid sink ID": {
			dataset: datasetInvalidSinkID,
			token:   token,
			err:     policies.ErrMalformedEntity,
		},
		"validate a dataset with a invalid policy ID": {
			dataset: datasetInvalidPolicyID,
			token:   token,
			err:     policies.ErrMalformedEntity,
		},
		"validate a dataset with a invalid agent group ID": {
			dataset: datasetInvalidAgentGroupID,
			token:   token,
			err:     policies.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.ValidateDataset(context.Background(), tc.token, tc.dataset)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestRetrieveDatasetByID(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	dataset := createDataset(t, svc, "dataset")

	cases := map[string]struct {
		id    string
		token string
		err   error
	}{
		"view an existing dataset": {
			id:    dataset.ID,
			token: token,
			err:   nil,
		},
		"view policy with wrong credentials": {
			id:    dataset.ID,
			token: "wrong",
			err:   policies.ErrUnauthorizedAccess,
		},
		"view non-existing policy": {
			id:    "9bb1b244-a199-93c2-aa03-28067b431e2c",
			token: token,
			err:   policies.ErrNotFound,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.ViewDatasetByID(context.Background(), tc.token, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestListDataset(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	var datasetList []policies.Dataset
	for i := 0; i < limit; i++ {
		pl := createDataset(t, svc, fmt.Sprintf("dataset-%d", i))
		datasetList = append(datasetList, pl)
	}

	cases := map[string]struct {
		token string
		pm    policies.PageMetadata
		size  uint64
		err   error
	}{
		"retrieve a list of datasets": {
			token: token,
			pm: policies.PageMetadata{
				Limit:  limit,
				Offset: 0,
			},
			size: limit,
			err:  nil,
		},
		"list half": {
			token: token,
			pm: policies.PageMetadata{
				Offset: limit / 2,
				Limit:  limit,
			},
			size: limit / 2,
			err:  nil,
		},
		"list last dataset": {
			token: token,
			pm: policies.PageMetadata{
				Offset: limit - 1,
				Limit:  limit,
			},
			size: 1,
			err:  nil,
		},
		"list empty set": {
			token: token,
			pm: policies.PageMetadata{
				Offset: limit + 1,
				Limit:  limit,
			},
			size: 0,
			err:  nil,
		},
		"list with zero limit": {
			token: token,
			pm: policies.PageMetadata{
				Offset: 1,
				Limit:  0,
			},
			size: 0,
			err:  nil,
		},
		"list with wrong credentials": {
			token: "wrong",
			pm: policies.PageMetadata{
				Offset: 0,
				Limit:  0,
			},
			size: 0,
			err:  policies.ErrUnauthorizedAccess,
		},
		"list all datasets sorted by name ascendant": {
			token: token,
			pm: policies.PageMetadata{
				Offset: 0,
				Limit:  limit,
				Order:  "name",
				Dir:    "asc",
			},
			size: limit,
			err:  nil,
		},
		"list all dataset sorted by name descendent": {
			token: token,
			pm: policies.PageMetadata{
				Offset: 0,
				Limit:  limit,
				Order:  "name",
				Dir:    "desc",
			},
			size: limit,
			err:  nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			page, err := svc.ListDatasets(context.Background(), tc.token, tc.pm)
			size := uint64(len(page.Datasets))
			assert.Equal(t, size, tc.size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, size))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
			testSortDataset(t, tc.pm, page.Datasets)
		})

	}
}

func TestListPoliciesByGroupIDInternal(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	agentGroupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := createPolicy(t, svc, "policy")

	var total = 10

	datasetsID := make([]string, total)

	for i := 0; i < total; i++ {
		ID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		sinkIDs := make([]string, 2)
		for i := 0; i < 2; i++ {
			sinkID, err := uuid.NewV4()
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
			sinkIDs = append(sinkIDs, sinkID.String())
		}

		validName, err := types.NewIdentifier(fmt.Sprintf("dataset-%d", i))
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		dataset := policies.Dataset{
			ID:           ID.String(),
			Name:         validName,
			PolicyID:     policy.ID,
			AgentGroupID: agentGroupID.String(),
			SinkIDs:      sinkIDs,
		}

		ds, err := svc.AddDataset(context.Background(), token, dataset)
		if err != nil {
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
		}

		datasetsID[i] = ds.ID
	}
	oID, _ := identify(token, users)

	listPlTest := make([]policies.PolicyInDataset, total)
	for i := 0; i < total; i++ {
		listPlTest[i] = policies.PolicyInDataset{
			Policy:    policy,
			DatasetID: datasetsID[i],
		}
	}

	cases := map[string]struct {
		ownerID  string
		groupID  []string
		policies []policies.PolicyInDataset
		size     uint64
		err      error
	}{
		"retrieve a list of policies by groupID": {
			ownerID:  oID,
			groupID:  []string{agentGroupID.String()},
			policies: listPlTest,
			size:     uint64(total),
			err:      nil,
		},
		"retrieve a list of policies by non-existent groupID": {
			ownerID:  oID,
			groupID:  []string{oID},
			policies: []policies.PolicyInDataset{},
			size:     uint64(0),
			err:      nil,
		},
		"list with empty ownerID": {
			ownerID:  "",
			groupID:  []string{agentGroupID.String()},
			policies: nil,
			size:     0,
			err:      policies.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			policies, err := svc.ListPoliciesByGroupIDInternal(context.Background(), tc.groupID, tc.ownerID)
			size := uint64(len(policies))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, size))
			assert.Equal(t, tc.policies, policies, fmt.Sprintf("%s: expected %p got %p", desc, tc.policies, policies))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})

	}
}

func TestRetrievePolicyByIDInternal(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	policy := createPolicy(t, svc, "policy")

	oID, _ := identify(token, users)

	wrongPlID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	cases := map[string]struct {
		policyID string
		ownerID  string
		err      error
	}{
		"view a existing policy": {
			policyID: policy.ID,
			ownerID:  oID,
			err:      nil,
		},
		"view policy with empty ownerID": {
			policyID: policy.ID,
			ownerID:  "",
			err:      policies.ErrMalformedEntity,
		},
		"view non-existing policy": {
			policyID: wrongPlID.String(),
			ownerID:  oID,
			err:      policies.ErrNotFound,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.ViewPolicyByIDInternal(context.Background(), tc.policyID, tc.ownerID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestListDatasetsByPolicyIDInternal(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	wrongPlID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := createPolicy(t, svc, "policy")

	var total = 10

	datasetsTest := make([]policies.Dataset, total)

	for i := 0; i < total; i++ {
		ID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		agentGroupID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		sinkIDs := make([]string, 2)
		for i := 0; i < 2; i++ {
			sinkID, err := uuid.NewV4()
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
			sinkIDs = append(sinkIDs, sinkID.String())
		}

		validName, err := types.NewIdentifier(fmt.Sprintf("dataset-%d", i))
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		dataset := policies.Dataset{
			ID:           ID.String(),
			Name:         validName,
			PolicyID:     policy.ID,
			AgentGroupID: agentGroupID.String(),
			SinkIDs:      sinkIDs,
		}

		ds, err := svc.AddDataset(context.Background(), token, dataset)
		if err != nil {
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
		}

		datasetsTest[i] = ds
	}

	cases := map[string]struct {
		token    string
		policyId string
		size     uint64
		err      error
	}{
		"retrieve a list of datasets by policyID": {
			policyId: policy.ID,
			token:    token,
			size:     uint64(total),
			err:      nil,
		},
		"retrieve a list of datasets by non-existent policyID": {
			policyId: wrongPlID.String(),
			token:    token,
			size:     uint64(0),
			err:      nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			datasets, err := svc.ListDatasetsByPolicyIDInternal(context.Background(), tc.policyId, tc.token)
			size := uint64(len(datasets))
			assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, size))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})

	}
}

func TestRetrieveDatasetByIDInternal(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	dataset := createDataset(t, svc, "dataset")

	oID, _ := identify(token, users)

	wrongPlID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	cases := map[string]struct {
		datasetID string
		ownerID   string
		err       error
	}{
		"view a existing dataset": {
			datasetID: dataset.ID,
			ownerID:   oID,
			err:       nil,
		},
		"view non-existing policy": {
			datasetID: wrongPlID.String(),
			ownerID:   oID,
			err:       policies.ErrNotFound,
		},
	}
	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.ViewDatasetByIDInternal(context.Background(), tc.ownerID, tc.datasetID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestInactivateDatasetsByGroupID(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	wrongGroupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	agentGroupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := createPolicy(t, svc, "policy")

	var total = 10

	datasetsTest := make([]policies.Dataset, total)

	for i := 0; i < total; i++ {
		ID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		sinkIDs := make([]string, 2)
		for i := 0; i < 2; i++ {
			sinkID, err := uuid.NewV4()
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
			sinkIDs = append(sinkIDs, sinkID.String())
		}

		validName, err := types.NewIdentifier(fmt.Sprintf("dataset-%d", i))
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		dataset := policies.Dataset{
			ID:           ID.String(),
			Name:         validName,
			PolicyID:     policy.ID,
			AgentGroupID: agentGroupID.String(),
			SinkIDs:      sinkIDs,
		}

		ds, err := svc.AddDataset(context.Background(), token, dataset)
		if err != nil {
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
		}

		datasetsTest[i] = ds
	}

	cases := map[string]struct {
		token   string
		groupID string
		err     error
	}{
		"inactivate a set of datasets by groupID": {
			groupID: agentGroupID.String(),
			token:   token,
			err:     nil,
		},
		"inactivate datasets with a non-existent groupID": {
			groupID: wrongGroupID.String(),
			token:   token,
			err:     policies.ErrNotFound,
		},
		"inactivate datasets with empty groupID": {
			groupID: "",
			token:   token,
			err:     policies.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := svc.InactivateDatasetByGroupID(context.Background(), tc.groupID, tc.token)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})

	}
}

func createPolicy(t *testing.T, svc policies.Service, name string) policies.Policy {
	t.Helper()
	ID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	validName, err := types.NewIdentifier(name)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := policies.Policy{
		ID:         ID.String(),
		Name:       validName,
		Backend:    "pktvisor",
		Format:     format,
		PolicyData: policy_data,
	}

	res, err := svc.AddPolicy(context.Background(), token, policy)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
	return res
}

func createDataset(t *testing.T, svc policies.Service, name string) policies.Dataset {
	t.Helper()
	ID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	agentGroupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	sinkIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		sinkID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
		sinkIDs[i] = sinkID.String()
	}

	validName, err := types.NewIdentifier(name)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	dataset := policies.Dataset{
		ID:           ID.String(),
		Name:         validName,
		PolicyID:     policyID.String(),
		AgentGroupID: agentGroupID.String(),
		SinkIDs:      sinkIDs,
	}

	res, err := svc.AddDataset(context.Background(), token, dataset)
	if err != nil {
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
	}
	return res
}

func TestDeleteSinkFromDataset(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	agentGroupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := createPolicy(t, svc, "policy")

	var total = 10

	datasetsTest := make([]policies.Dataset, total)

	sinkIDs := make([]string, 0)
	for i := 0; i < total; i++ {
		ID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		for i := 0; i < 2; i++ {
			sinkID, err := uuid.NewV4()
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
			sinkIDs = append(sinkIDs, sinkID.String())
		}

		validName, err := types.NewIdentifier(fmt.Sprintf("dataset-%d", i))
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		dataset := policies.Dataset{
			ID:           ID.String(),
			Name:         validName,
			PolicyID:     policy.ID,
			AgentGroupID: agentGroupID.String(),
			SinkIDs:      sinkIDs,
		}

		ds, err := svc.AddDataset(context.Background(), token, dataset)
		if err != nil {
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
		}

		datasetsTest[i] = ds
	}

	cases := map[string]struct {
		ownerID string
		sinkID  string
		err     error
	}{
		"delete sinkID of a set of datasets": {
			sinkID:  sinkIDs[0],
			ownerID: datasetsTest[0].MFOwnerID,
			err:     nil,
		},
		"delete sinkID of a set of datasets with empty sinkID": {
			sinkID:  "",
			ownerID: datasetsTest[0].MFOwnerID,
			err:     policies.ErrMalformedEntity,
		},
		"delete sinkID of a set of datasets with empty owner": {
			sinkID:  sinkIDs[0],
			ownerID: "",
			err:     policies.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.DeleteSinkFromAllDatasetsInternal(context.Background(), tc.sinkID, tc.ownerID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestInactivateDatasetByID(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	agentGroupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := createPolicy(t, svc, "policy")

	var total = 10

	datasetsTest := make([]policies.Dataset, total)
	datasetsID := make([]string, 0, total)

	sinkIDs := make([]string, 0)
	for i := 0; i < total; i++ {
		ID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		datasetsID = append(datasetsID, ID.String())

		for i := 0; i < 2; i++ {
			sinkID, err := uuid.NewV4()
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
			sinkIDs = append(sinkIDs, sinkID.String())
		}

		validName, err := types.NewIdentifier(fmt.Sprintf("dataset-%d", i))
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		dataset := policies.Dataset{
			ID:           ID.String(),
			Name:         validName,
			PolicyID:     policy.ID,
			AgentGroupID: agentGroupID.String(),
			SinkIDs:      sinkIDs,
		}

		ds, err := svc.AddDataset(context.Background(), token, dataset)
		if err != nil {
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
		}

		datasetsTest[i] = ds
	}

	cases := map[string]struct {
		ownerID    string
		datasetIDs []string
		err        error
	}{
		"inactivate a set of datasets by ID": {
			datasetIDs: datasetsID,
			ownerID:    datasetsTest[0].MFOwnerID,
			err:        nil,
		},
		"inactivate datasets with empty ownerID": {
			datasetIDs: datasetsID,
			ownerID:    "",
			err:        policies.ErrMalformedEntity,
		},
		"inactivate datasets with empty ID": {
			datasetIDs: []string{""},
			ownerID:    datasetsTest[0].MFOwnerID,
			err:        policies.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			for _, id := range tc.datasetIDs {
				err := svc.InactivateDatasetByIDInternal(context.Background(), tc.ownerID, id)
				assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
			}
		})
	}
}

func TestDeleteAGroupFromDataset(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	agentGroupID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := createPolicy(t, svc, "policy")

	var total = 10

	datasetsTest := make([]policies.Dataset, total)

	sinkIDs := make([]string, 0)
	for i := 0; i < total; i++ {
		ID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		for i := 0; i < 2; i++ {
			sinkID, err := uuid.NewV4()
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
			sinkIDs = append(sinkIDs, sinkID.String())
		}

		validName, err := types.NewIdentifier(fmt.Sprintf("dataset-%d", i))
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

		dataset := policies.Dataset{
			ID:           ID.String(),
			Name:         validName,
			PolicyID:     policy.ID,
			AgentGroupID: agentGroupID.String(),
			SinkIDs:      sinkIDs,
		}

		ds, err := svc.AddDataset(context.Background(), token, dataset)
		if err != nil {
			require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
		}

		datasetsTest[i] = ds
	}

	cases := map[string]struct {
		token  string
		aGroup string
		err    error
	}{
		"delete agent group of a set of datasets": {
			aGroup: agentGroupID.String(),
			token:  token,
			err:    nil,
		},
		"delete agent group of a set of datasets with empty agent group ID": {
			aGroup: "",
			token:  token,
			err:    policies.ErrMalformedEntity,
		},
		"delete agent group of a set of datasets with empty owner": {
			aGroup: agentGroupID.String(),
			token:  "wrong",
			err:    policies.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := svc.DeleteAgentGroupFromAllDatasets(context.Background(), tc.aGroup, tc.token)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func testSortPolicies(t *testing.T, pm policies.PageMetadata, ags []policies.Policy) {
	t.Helper()
	switch pm.Order {
	case "name":
		current := ags[0]
		for _, res := range ags {
			if pm.Dir == "asc" {
				assert.GreaterOrEqual(t, res.Name.String(), current.Name.String())
			}
			if pm.Dir == "desc" {
				assert.GreaterOrEqual(t, current.Name.String(), res.Name.String())
			}
			current = res
		}
	default:
		break
	}
}

func TestDuplicatePolicy(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	policySingleCopy := createPolicy(t, svc, "policy-single")

	policyMultipleCopy := createPolicy(t, svc, "policy")
	copy := createPolicy(t, svc, "policy_copy")
	for i := 2; i <= 10; i++ {
		name := fmt.Sprintf("policy_copy%d", i)
		_ = createPolicy(t, svc, name)
	}

	cases := map[string]struct {
		policyID     string
		policyName   string
		expectedName string
		token        string
		err          error
	}{
		"duplicate a existing policy with no name specified": {
			policyID:     policySingleCopy.ID,
			policyName:   "",
			expectedName: "policy-single_copy",
			token:        token,
			err:          nil,
		},
		"duplicate a existing policy with a specified name": {
			policyID:     policySingleCopy.ID,
			policyName:   "policyDuplicated",
			expectedName: "policyDuplicated",
			token:        token,
			err:          nil,
		},
		"duplicate a existing policy with an invalid token": {
			policyID: policySingleCopy.ID,
			token:    invalidToken,
			err:      policies.ErrUnauthorizedAccess,
		},
		"duplicate a existing policy with empty ID": {
			policyID: "",
			token:    token,
			err:      policies.ErrNotFound,
		},
		"duplicate a existing policy with more than three copies": {
			policyID:     policyMultipleCopy.ID,
			expectedName: "",
			token:        token,
			err:          errors.ErrConflict,
		},
		"duplicate a copy of a policy": {
			policyID:     copy.ID,
			expectedName: "policy_copy_copy",
			token:        token,
			err:          nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			policyDuplicated, err := svc.DuplicatePolicy(context.Background(), tc.token, tc.policyID, tc.policyName)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, err, tc.err))

			originalPolicy, _ := svc.ViewPolicyByID(context.Background(), token, tc.policyID)
			if err == nil {
				assert.Equal(t, originalPolicy.Policy, policyDuplicated.Policy, fmt.Sprintf("%s: expected %v got %v", desc, originalPolicy.Policy, policyDuplicated.Policy))
				assert.Equal(t, originalPolicy.PolicyData, policyDuplicated.PolicyData, fmt.Sprintf("%s: expected %v got %v", desc, originalPolicy.PolicyData, policyDuplicated.PolicyData))
				assert.Equal(t, originalPolicy.Format, policyDuplicated.Format, fmt.Sprintf("%s: expected %v got %v", desc, originalPolicy.Format, policyDuplicated.Format))
				assert.Equal(t, originalPolicy.Backend, policyDuplicated.Backend, fmt.Sprintf("%s: expected %v got %v", desc, originalPolicy.Backend, policyDuplicated.Backend))

				assert.Equal(t, tc.expectedName, policyDuplicated.Name.String(), fmt.Sprintf("%s: expected %v got %v", desc, tc.expectedName, policyDuplicated.Name))
			}
		})
	}
}

func testSortDataset(t *testing.T, pm policies.PageMetadata, ags []policies.Dataset) {
	t.Helper()
	switch pm.Order {
	case "name":
		current := ags[0]
		for _, res := range ags {
			if pm.Dir == "asc" {
				assert.GreaterOrEqual(t, res.Name.String(), current.Name.String())
			}
			if pm.Dir == "desc" {
				assert.GreaterOrEqual(t, current.Name.String(), res.Name.String())
			}
			current = res
		}
	default:
		break
	}
}

func identify(token string, auth mainflux.AuthServiceClient) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return res.GetId(), nil
}
