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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	token       = "token"
	invalidToken = "invalid"
	email       = "user@example.com"
	format      = "yaml"
	policy_data = `version: "1.0"
visor:
  taps:
    anycast:
      type: pcap
      config:
        iface: eth0`
	limit = 10
)

func newService(auth mainflux.AuthServiceClient) policies.Service {
	policyRepo := plmocks.NewPoliciesRepository()
	return policies.New(nil, auth, policyRepo, nil)
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
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			res, err := svc.EditPolicy(context.Background(), tc.token, tc.pol, tc.format, tc.policyData)
			assert.Equal(t, tc.pol.Name.String(), res.Name.String(), fmt.Sprintf("%s: expected name %s got %s", desc, tc.pol.Name.String(), res.Name.String()))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected error %d got %d", desc, tc.err, err))
		})
	}

}

func TestValidatePolicy(t *testing.T) {
	var nameID, _ = types.NewIdentifier("my-policy")
	var policy = policies.Policy{
		Name: nameID,
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

func createPolicy(t *testing.T, svc policies.Service, name string) policies.Policy {
	t.Helper()
	ID, err := uuid.NewV4()
	if err != nil {
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
	}

	validName, err := types.NewIdentifier(name)
	if err != nil {
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
	}

	policy := policies.Policy{
		ID:      ID.String(),
		Name:    validName,
		Backend: "pktvisor",
	}

	res, err := svc.AddPolicy(context.Background(), token, policy, format, policy_data)
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
