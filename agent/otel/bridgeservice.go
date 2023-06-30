package otel

import (
	"context"
	"github.com/orb-community/orb/agent/policies"
	"strings"
)

type AgentBridgeService interface {
	RetrieveAgentInfoByPolicyName(policyName string) (*AgentDataPerPolicy, error)
	NotifyAgentDisconnection(ctx context.Context, err error)
}

type AgentDataPerPolicy struct {
	AgentName string
	PolicyID  string
	Datasets  string
	AgentTags map[string]string
}

var _ AgentBridgeService = (*BridgeService)(nil)

type BridgeService struct {
	bridgeContext context.Context
	AgentName     string
	policyRepo    policies.PolicyRepo
	AgentTags     map[string]string
}

func NewBridgeService(ctx context.Context, agentName string, policyRepo *policies.PolicyRepo, agentTags map[string]string) *BridgeService {
	return &BridgeService{
		AgentName:     agentName,
		bridgeContext: ctx,
		policyRepo:    *policyRepo,
		AgentTags:     agentTags,
	}
}

func (b *BridgeService) RetrieveAgentInfoByPolicyName(policyName string) (*AgentDataPerPolicy, error) {
	pData, err := b.policyRepo.GetByName(policyName)
	if err != nil {
		return nil, err
	}
	return &AgentDataPerPolicy{
		AgentName: b.AgentName,
		PolicyID:  pData.ID,
		Datasets:  strings.Join(pData.GetDatasetIDs(), ","),
		AgentTags: b.AgentTags,
	}, nil
}

func (b *BridgeService) NotifyAgentDisconnection(ctx context.Context, err error) {
	ctx.Done()
	b.bridgeContext.Done()
}
