package grpc

type agentRes struct {
	id      string
	name    string
	channel string
}

type emptyRes struct {
	err error
}
