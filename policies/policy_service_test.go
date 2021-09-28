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
)

const (
	token        = "token"
	invalidToken = "invalid"
	email        = "user@example.com"
	format       = "yaml"
	policy_data  = `version: "1.0"
visor:
  taps:
    anycast:
      type: pcap
      config:
        iface: eth0`
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
		ID:        policy.ID,
		Name:      nameID,
		MFOwnerID: policy.MFOwnerID,
	}

	cases := map[string]struct {
		pol        policies.Policy
		token      string
		format     string
		policyData string
		err        error
	}{
		"update a existing policy": {
			pol:        newPolicy,
			token:      token,
			format:     format,
			policyData: policy_data,
			err:        nil,
		},
		"update policy with wrong credentials": {
			pol:        newPolicy,
			token:      "invalidToken",
			format:     format,
			policyData: policy_data,
			err:        policies.ErrUnauthorizedAccess,
		},
		"update a non-existing policy": {
			pol:        wrongPolicy,
			token:      token,
			format:     format,
			policyData: policy_data,
			err:        policies.ErrNotFound,
		},
		"update a existing policy with invalid format": {
			pol:        newPolicy,
			token:      token,
			format:     "invalid",
			policyData: policy_data,
			err:        policies.ErrValidatePolicy,
		},
		"update a existing policy with invalid policy_data": {
			pol:        newPolicy,
			token:      token,
			format:     format,
			policyData: "invalid",
			err:        policies.ErrValidatePolicy,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			res, err := svc.EditPolicy(context.Background(), tc.token, tc.pol, tc.format, tc.policyData)
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
		Name:    nameID,
		Backend: "pktvisor",
		OrbTags: map[string]string{"region": "eu", "node_type": "dns"},
	}

	users := flmocks.NewAuthService(map[string]string{token: email})
	svc := newService(users)

	cases := map[string]struct {
		policy     policies.Policy
		token      string
		format     string
		policyData string
		err        error
	}{
		"validate a new policy": {
			policy:     policy,
			token:      token,
			format:     format,
			policyData: policy_data,
			err:        nil,
		},
		"validate a policy with a invalid token": {
			policy:     policy,
			token:      invalidToken,
			format:     format,
			policyData: policy_data,
			err:        policies.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := svc.ValidatePolicy(context.Background(), tc.token, policy, format, policy_data)
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
	}

	cases := map[string]struct {
		policy policies.Policy
		format string
		token  string
		err    error
	}{
		"create a new policy": {
			policy: policy,
			format: format,
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
			_, err := svc.AddPolicy(context.Background(), tc.token, tc.policy, tc.format, policy_data)
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

	policy := createDataset(t, svc, "policy")

	nameID, err := types.NewIdentifier("new-policy")
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	wrongOwnerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	wrongDataset := policies.Dataset{MFOwnerID: wrongOwnerID.String()}
	newDataset := policies.Dataset{
		ID:        policy.ID,
		Name:      nameID,
		MFOwnerID: policy.MFOwnerID,
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
			err:   policies.ErrNotFound,
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
		sinkIDsArray = []string{"f5b2d342-211d-a9ab-1233-63199a3fc16f", "03679425-aa69-4574-bf62-e0fe71b80939"}
		dataset                    = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: policy.ID, SinkIDs: sinkIDsArray, Valid: true}
		datasetEmptySinkID         = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: policy.ID, SinkIDs: []string{}, Valid: true}
		datasetEmptyPolicyID       = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: "", SinkIDs: sinkIDsArray, Valid: true}
		datasetEmptyAgentGroupID   = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "", PolicyID: policy.ID, SinkIDs: sinkIDsArray, Valid: true}
		datasetInvalidSinkID       = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: policy.ID, SinkIDs: []string{"invalid"}, Valid: true}
		datasetInvalidPolicyID     = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "8fd6d12d-6a26-5d85-dc35-f9ba8f4d93db", PolicyID: "invalid", SinkIDs: sinkIDsArray, Valid: true}
		datasetInvalidAgentGroupID = policies.Dataset{Name: nameID, Tags: map[string]string{"region": "eu", "node_type": "dns"}, AgentGroupID: "invalid", PolicyID: policy.ID, SinkIDs: sinkIDsArray, Valid: true}
	)

	cases := map[string]struct {
		dataset  policies.Dataset
		token    string
		err      error
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
			page, err := svc.ListDataset(context.Background(), tc.token, tc.pm)
			size := uint64(len(page.Datasets))
			assert.Equal(t, size, tc.size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, size))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
			testSortDataset(t, tc.pm, page.Datasets)
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
		ID:      ID.String(),
		Name:    validName,
		Backend: "pktvisor",
	}

	res, err := svc.AddPolicy(context.Background(), token, policy, format, policy_data)
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
		sinkIDs = append(sinkIDs, sinkID.String())
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
