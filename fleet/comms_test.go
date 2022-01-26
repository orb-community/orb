package fleet_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mainflux/mainflux"
	mflog "github.com/mainflux/mainflux/logger"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/ns1labs/orb/fleet"
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
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"os"
	"testing"
)

const bufSize = 1024 * 1024

var (
	lis *bufconn.Listener
	users = flmocks.NewAuthService(map[string]string{token: email})
	policiesSVC = newPoliciesService(users)
)

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


func newCommsService() fleet.AgentCommsService {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	agentGroupRepo := flmocks.NewAgentGroupRepository()
	agentRepo := flmocks.NewAgentRepositoryMock()

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
	agentPubSub, err := mfnats.NewPubSub(url.URL, "fleet", mflogger)
	if err != nil {
		log.Fatalf("Failed to create PubSub %v", err)
	}

	return fleet.NewFleetCommsService(logger, policyClient, agentRepo, agentGroupRepo, agentPubSub)
}

func TestNotifyGroupNewDataset(t *testing.T){
	commsSVC := newCommsService()

	thingsServer := newThingsServer(newThingsService(users))
	fleetSVC := newService(users, thingsServer.URL)

	ag, err := createAgentGroup(t, "group", fleetSVC)
	if err != nil {
		log.Fatalf("Faled to create agent group: %v", err)
	}

	policy := createPolicy(t, policiesSVC, "policy")
	dataset := createDataset(t, policiesSVC, "dataset")

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
	}

	for desc, tc := range cases{
		err := commsSVC.NotifyGroupNewDataset(context.Background(), tc.agentGroup, tc.datasetID, tc.policyID, tc.ownerID)
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
