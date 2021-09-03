// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	thmocks "github.com/mainflux/mainflux/things/mocks"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/policies"
	"github.com/ns1labs/orb/policies/mocks"
	"github.com/ns1labs/orb/policies/pb"
	"net"
	"os"
	"testing"

	policiesgrpc "github.com/ns1labs/orb/policies/api/grpc"
	"github.com/opentracing/opentracing-go/mocktracer"
	"google.golang.org/grpc"
)

const (
	port  = 18080
	token = "token"
	email = "john.doe@email.com"
)

var (
	svc     policies.Service
	policy  policies.Policy
	dataset policies.Dataset
)

func TestMain(m *testing.M) {
	startServer()
	code := m.Run()
	os.Exit(code)
}

func startServer() {
	svc = newService(map[string]string{token: email})
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	server := grpc.NewServer()
	pb.RegisterPolicyServiceServer(server, policiesgrpc.NewServer(mocktracer.New(), svc))
	go server.Serve(listener)
}

func newService(tokens map[string]string) policies.Service {
	auth := thmocks.NewAuthService(tokens)
	repo := mocks.NewPoliciesRepository()

	oID, _ := uuid.NewV4()
	pname, _ := types.NewIdentifier("testpolicy")

	policy = policies.Policy{
		Name:      pname,
		MFOwnerID: oID.String(),
	}
	policyid, _ := repo.SavePolicy(context.Background(), policy)
	policy.ID = policyid

	gID, _ := uuid.NewV4()
	gname, _ := types.NewIdentifier("testdataset")
	dataset = policies.Dataset{
		Name:         gname,
		MFOwnerID:    oID.String(),
		AgentGroupID: gID.String(),
		PolicyID:     policyid,
	}
	datasetid, _ := repo.SaveDataset(context.Background(), dataset)
	dataset.ID = datasetid

	return policies.New(nil, auth, repo, nil)
}
