package fleet_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mainflux/mainflux"
	mflog "github.com/mainflux/mainflux/logger"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/backend/pktvisor"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/ns1labs/orb/pkg/config"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
	policyGRPC "github.com/ns1labs/orb/policies/api/grpc"
	plmocks "github.com/ns1labs/orb/policies/mocks"
	"github.com/ns1labs/orb/policies/pb"
	sinkmocks "github.com/ns1labs/orb/sinks/mocks"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"os"
	"testing"
)

const bufSize = 1024 * 1024

var (
	lis         *bufconn.Listener
	users       = flmocks.NewAuthService(map[string]string{token: email})
	policiesSVC = newPoliciesService(users)
)

func newFleetService(auth mainflux.AuthServiceClient, url string, agentGroupRepo fleet.AgentGroupRepository, agentRepo fleet.AgentRepository) fleet.Service {
	agentComms := flmocks.NewFleetCommService(agentRepo, agentGroupRepo)
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("%v", err)
	}
	config := mfsdk.Config{
		BaseURL: url,
	}

	mfsdk := mfsdk.NewSDK(config)
	pktvisor.Register(auth, agentRepo)
	return fleet.NewFleetService(logger, auth, agentRepo, agentGroupRepo, agentComms, mfsdk)
}

func newPoliciesService(auth mainflux.AuthServiceClient) policies.Service {
	policyRepo := plmocks.NewPoliciesRepository()
	fleetGrpcClient := flmocks.NewClient()
	SinkServiceClient := sinkmocks.NewClient()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("%v", err)
	}

	return policies.New(logger, auth, policyRepo, fleetGrpcClient, SinkServiceClient)
}

func init() {
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()

	tracer := mocktracer.New()
	policyServer := policyGRPC.NewServer(tracer, policiesSVC)

	pb.RegisterPolicyServiceServer(server, policyServer)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func newCommsService(agentGroupRepo fleet.AgentGroupRepository, agentRepo fleet.AgentRepository) fleet.AgentCommsService {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}
	policyClient := pb.NewPolicyServiceClient(conn)

	mflogger, err := mflog.New(os.Stdout, "debug")
	if err != nil {
		log.Fatalf(err.Error())
	}

	url := config.LoadNatsConfig("orb_fleet")
	agentPubSub, err := flmocks.NewPubSub(url.URL, "fleet", mflogger)
	if err != nil {
		log.Fatalf("Failed to create PubSub %v", err)
	}

	return fleet.NewFleetCommsService(logger, policyClient, agentRepo, agentGroupRepo, agentPubSub)
}

func TestNotifyGroupNewDataset(t *testing.T) {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

	commsSVC := newCommsService(agentGroupRepo, agentRepo)

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newFleetService(users, thingsServer.URL, agentGroupRepo, agentRepo)

	ag, err := createAgentGroup(t, "group", fleetSVC)
	assert.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policy := createPolicy(t, policiesSVC, "policy")
	dataset := createDataset(t, policiesSVC, "dataset", ag.ID)

	wrongPolicyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		policyID   string
		ownerID    string
		datasetID  string
		agentGroup fleet.AgentGroup
		err        error
	}{
		"Notify a existent group new dataset": {
			ownerID:    ag.MFOwnerID,
			policyID:   policy.ID,
			datasetID:  dataset.ID,
			agentGroup: ag,
			err:        nil,
		},
		"Notify a existent group new dataset with wrong policyID": {
			ownerID:    ag.MFOwnerID,
			policyID:   wrongPolicyID.String(),
			datasetID:  dataset.ID,
			agentGroup: ag,
			err:        status.Error(codes.Internal, "internal server error"),
		},
	}

	for desc, tc := range cases {
		err := commsSVC.NotifyGroupNewDataset(context.Background(), tc.agentGroup, tc.datasetID, tc.policyID, tc.ownerID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestNotifyGroupPolicyRemoval(t *testing.T) {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

	commsSVC := newCommsService(agentGroupRepo, agentRepo)

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newFleetService(users, thingsServer.URL, agentGroupRepo, agentRepo)

	agent, err := createAgent(t, "agent", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	group, err := createAgentGroup(t, "group2", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		agent      fleet.Agent
		agentGroup fleet.AgentGroup
		policyID   string
		policyName string
		backend    string
		err        error
	}{
		"Notify group policy deletion": {
			agent:      agent,
			agentGroup: group,
			policyID:   policyID.String(),
			policyName: "policy2",
			backend:    "pktvisor",
			err:        nil,
		},
	}

	for desc, tc := range cases {
		err := commsSVC.NotifyGroupPolicyRemoval(tc.agentGroup, tc.policyID, tc.policyName, tc.backend)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestNotifyAgentAllDatasets(t *testing.T) {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

	commsSVC := newCommsService(agentGroupRepo, agentRepo)

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newFleetService(users, thingsServer.URL, agentGroupRepo, agentRepo)

	validAgentName, err := types.NewIdentifier("agent2")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	ag, err := fleetSVC.CreateAgent(context.Background(), "token", fleet.Agent{
		Name:      validAgentName,
		AgentTags: map[string]string{"test": "true"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	validGroupName, err := types.NewIdentifier("group3")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	group, err := fleetSVC.CreateAgentGroup(context.Background(), "token", fleet.AgentGroup{
		Name: validGroupName,
		Tags: map[string]string{"test": "true"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	_ = createDataset(t, policiesSVC, "dataset2", group.ID)

	noMatchingGroup, err := createAgent(t, "agent3", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		agent fleet.Agent
		err   error
	}{
		"Notify agent all policies should run": {
			agent: ag,
			err:   nil,
		},
		"Notify agent all policies with malformed agent": {
			agent: fleet.Agent{MFThingID: ""},
			err:   errors.ErrMalformedEntity,
		},
		"Notify agent with no matching groups": {
			agent: noMatchingGroup,
			err:   nil,
		},
		"Notify agent with wrong thingID": {
			agent: fleet.Agent{
				MFOwnerID:   ag.MFOwnerID,
				MFThingID:   wrongID,
				MFChannelID: ag.MFChannelID,
				AgentTags:   ag.AgentTags,
			},
			err: fleet.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		err := commsSVC.NotifyAgentAllDatasets(tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestNotifyGroupRemoval(t *testing.T) {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

	commsSVC := newCommsService(agentGroupRepo, agentRepo)

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newFleetService(users, thingsServer.URL, agentGroupRepo, agentRepo)

	group, err := createAgentGroup(t, "group4", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		agentGroup fleet.AgentGroup
		err        error
	}{
		"Notify group deletion": {
			agentGroup: group,
			err:        nil,
		},
	}

	for desc, tc := range cases {
		err := commsSVC.NotifyGroupRemoval(tc.agentGroup)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestNotifyGroupPolicyUpdate(t *testing.T) {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

	commsSVC := newCommsService(agentGroupRepo, agentRepo)

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newFleetService(users, thingsServer.URL, agentGroupRepo, agentRepo)

	policy := createPolicy(t, policiesSVC, "policy3")

	agent, err := createAgent(t, "agent4", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	group, err := createAgentGroup(t, "group6", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		agent      fleet.Agent
		agentGroup fleet.AgentGroup
		policyID   string
		ownerID    string
		err        error
	}{
		"Notify group a policy update": {
			agent:      agent,
			agentGroup: group,
			policyID:   policy.ID,
			ownerID:    policy.MFOwnerID,
			err:        nil,
		},
		"Notify group a policy update wih wrong policyID": {
			agent:      agent,
			agentGroup: group,
			policyID:   wrongID,
			ownerID:    policy.MFOwnerID,
			err:        status.Error(codes.Internal, "internal server error"),
		},
	}

	for desc, tc := range cases {
		err := commsSVC.NotifyGroupPolicyUpdate(context.Background(), tc.agentGroup, tc.policyID, tc.ownerID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestNotifyAgentGroupMembership(t *testing.T) {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

	commsSVC := newCommsService(agentGroupRepo, agentRepo)

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newFleetService(users, thingsServer.URL, agentGroupRepo, agentRepo)

	validAgentName, err := types.NewIdentifier("agent5")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	ag, err := fleetSVC.CreateAgent(context.Background(), "token", fleet.Agent{
		Name:      validAgentName,
		AgentTags: map[string]string{"test": "true"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	validGroupName, err := types.NewIdentifier("group5")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	_, err = fleetSVC.CreateAgentGroup(context.Background(), "token", fleet.AgentGroup{
		Name: validGroupName,
		Tags: map[string]string{"test": "true"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	noMatchingGroup, err := createAgent(t, "agent6", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		agent fleet.Agent
		err   error
	}{
		"Notify agent all AgentGroup memberships it belongs to": {
			agent: ag,
			err:   nil,
		},
		"Notify agent not belong to any AgentGroup": {
			agent: noMatchingGroup,
			err:   nil,
		},
		"Notify agent, but missing thingID": {
			agent: fleet.Agent{
				MFOwnerID:   ag.MFOwnerID,
				MFThingID:   "",
				MFChannelID: ag.MFChannelID,
				AgentTags:   ag.AgentTags,
			},
			err: fleet.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		err := commsSVC.NotifyAgentGroupMemberships(tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestNotifyGroupDatasetRemoval(t *testing.T) {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

	commsSVC := newCommsService(agentGroupRepo, agentRepo)

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newFleetService(users, thingsServer.URL, agentGroupRepo, agentRepo)

	agent, err := createAgent(t, "agent4", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	group, err := createAgentGroup(t, "group6", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	policyID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	datasetID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		agent      fleet.Agent
		agentGroup fleet.AgentGroup
		dsID       string
		policyID   string
		err        error
	}{
		"Notify group dataset deletion": {
			agent:      agent,
			agentGroup: group,
			dsID:       datasetID.String(),
			policyID:   policyID.String(),
			err:        nil,
		},
	}

	for desc, tc := range cases {
		err := commsSVC.NotifyGroupDatasetRemoval(tc.agentGroup, tc.dsID, tc.policyID)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestNotifyAgentStop(t *testing.T) {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

	commsSVC := newCommsService(agentGroupRepo, agentRepo)

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newFleetService(users, thingsServer.URL, agentGroupRepo, agentRepo)

	agent, err := createAgent(t, "agent4", fleetSVC)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		agent  fleet.Agent
		reason string
		err    error
	}{
		"Notify agent to stop": {
			agent:  agent,
			reason: "",
			err:    nil,
		},
	}

	for desc, tc := range cases {
		err := commsSVC.NotifyAgentStop(tc.agent, tc.reason)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func TestNotifyAgentNewGroupMembership(t *testing.T) {
	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

	commsSVC := newCommsService(agentGroupRepo, agentRepo)

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newFleetService(users, thingsServer.URL, agentGroupRepo, agentRepo)

	validAgentName, err := types.NewIdentifier("agent5")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	ag, err := fleetSVC.CreateAgent(context.Background(), "token", fleet.Agent{
		Name:      validAgentName,
		AgentTags: map[string]string{"test": "true"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	validGroupName, err := types.NewIdentifier("group5")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	_, err = fleetSVC.CreateAgentGroup(context.Background(), "token", fleet.AgentGroup{
		Name: validGroupName,
		Tags: map[string]string{"test": "true"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		agent      fleet.Agent
		agentGroup fleet.AgentGroup
		err        error
	}{
		"Notify agent a new membership of AgentGroup": {
			agent:      ag,
			agentGroup: agentGroup,
			err:        nil,
		},
	}

	for desc, tc := range cases {
		err := commsSVC.NotifyAgentNewGroupMembership(tc.agent, tc.agentGroup)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
	}
}

func createPolicy(t *testing.T, svc policies.Service, name string) policies.Policy {
	t.Helper()
	ID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	validName, err := types.NewIdentifier(name)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policy := policies.Policy{
		Name:        validName,
		MFOwnerID:   ID.String(),
		Description: "An example policy",
		Backend:     "pktvisor",
		Version:     0,
		OrbTags:     map[string]string{"region": "eu"},
	}
	policy_data := `version: "1.0"
visor:
  taps:
    anycast:
      type: pcap
      config:
        iface: eth0`

	res, err := svc.AddPolicy(context.Background(), token, policy, "yaml", policy_data)
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
	return res
}

func createDataset(t *testing.T, svc policies.Service, name string, groupID string) policies.Dataset {
	t.Helper()
	ID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))

	policyID, err := uuid.NewV4()
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
		AgentGroupID: groupID,
		SinkIDs:      sinkIDs,
	}

	res, err := svc.AddDataset(context.Background(), token, dataset)
	if err != nil {
		require.Nil(t, err, fmt.Sprintf("Unexpected error: %s", err))
	}
	return res
}
