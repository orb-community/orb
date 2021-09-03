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
)

func newService(auth mainflux.AuthServiceClient) policies.Service {
	policyRepo := plmocks.NewPoliciesRepository()
	return policies.New(auth, policyRepo)
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
