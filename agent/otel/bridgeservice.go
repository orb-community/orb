package otel

import (
	"github.com/ns1labs/orb/agent/policies"
	"strings"
)

type AgentBridgeService interface {
	RetrieveAgentInfoByPolicyName(policyName string) (*AgentDataPerPolicy, error)
}

type AgentDataPerPolicy struct {
	PolicyID  string
	Datasets  string
	AgentTags string
}

var _ AgentBridgeService = (*bridgeService)(nil)

type bridgeService struct {
	policyRepo policies.PolicyRepo
	AgentTags  string
}

func NewBridgeService(policyRepo *policies.PolicyRepo, agentTags string) *bridgeService {
	return &bridgeService{
		policyRepo: *policyRepo,
		AgentTags:  agentTags,
	}
}

func (b *bridgeService) RetrieveAgentInfoByPolicyName(policyName string) (*AgentDataPerPolicy, error) {
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
