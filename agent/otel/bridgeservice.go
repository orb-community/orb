package otel

type AgentBridgeService interface {
	RetrieveAgentInfoByPolicyID(policyID string) (*AgentDataPerPolicy, error)
}

type AgentDataPerPolicy struct {
	PolicyID  string
	Datasets  string
	OrbTags   string
	AgentTags string
}

var _ AgentBridgeService = (*bridgeService)(nil)

type bridgeService struct {
}

func NewBridgeService() *bridgeService {
	return &bridgeService{}
}

func (b bridgeService) RetrieveAgentInfoByPolicyID(policyID string) (*AgentDataPerPolicy, error) {

	return nil, nil
}
