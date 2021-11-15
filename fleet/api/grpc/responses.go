package grpc

type agentRes struct {
	id      string
	name    string
	channel string
}

type agentGroupRes struct {
	id      string
	name    string
	channel string
}

type ownerRes struct {
	ownerID   string
	agentName string
}

type emptyRes struct {
	err error
}
