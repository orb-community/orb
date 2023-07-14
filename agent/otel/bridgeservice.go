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
	PolicyID  string
	Datasets  string
	AgentTags map[string]string
}

var _ AgentBridgeService = (*BridgeService)(nil)

type BridgeService struct {
	bridgeContext context.Context
	cancelFunc    context.CancelFunc
	policyRepo    policies.PolicyRepo
	AgentTags     map[string]string
}

func NewBridgeService(ctx context.Context, cancelFunc context.CancelFunc, policyRepo *policies.PolicyRepo, agentTags map[string]string) *BridgeService {
	return &BridgeService{
		bridgeContext: ctx,
		cancelFunc:    cancelFunc,
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
		PolicyID:  pData.ID,
		Datasets:  strings.Join(pData.GetDatasetIDs(), ","),
		AgentTags: b.AgentTags,
	}, nil
}

func (b *BridgeService) NotifyAgentDisconnection(ctx context.Context, err error) {
	ctx.Done()
	b.cancelFunc()
}
