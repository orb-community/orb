package migrate

import (
	"context"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks"
	"github.com/orb-community/orb/sinks/authentication_type"
	"go.uber.org/zap"
)

type Plan interface {
	Up(ctx context.Context) error
}

type Plan1UpdateConfiguration struct {
	logger          *zap.Logger
	service         sinks.SinkService
	passwordService authentication_type.PasswordService
}

func (p *Plan1UpdateConfiguration) Up(ctx context.Context) (err error) {
	allSinks, err := p.service.ListSinksInternal(ctx, sinks.Filter{})
	for _, sink := range allSinks {
		if _, ok := sink.Config["authentication"]; !ok {
			sinkRemoteHost, ok := sink.Config["remote_host"]
			if !ok {
				p.logger.Error("failed to update sink for lack of remote_host", zap.String("sinkID", sink.ID))
				continue
			}
			sinkUsername, ok := sink.Config["username"]
			if !ok {
				p.logger.Error("failed to update sink for lack of username", zap.String("sinkID", sink.ID))
				continue
			}
			encodedPassword, ok := sink.Config["password"]
			if !ok {
				p.logger.Error("failed to update sink for lack of password", zap.String("sinkID", sink.ID))
				continue
			}
			decodedPassword, err := p.passwordService.DecodePassword(encodedPassword.(string))
			if err != nil {
				p.logger.Error("failed to update sink for failure in decoding password",
					zap.String("sinkID", sink.ID), zap.Error(err))
				continue
			}
			newMetadata := types.Metadata{
				"authentication": types.Metadata{
					"type":     "basicauth",
					"username": sinkUsername,
					"password": decodedPassword,
				},
				"exporter": types.Metadata{
					"remote_host": sinkRemoteHost,
				},
			}
			sink.Config = newMetadata
			_, err = p.service.UpdateSinkInternal(ctx, sink)
			if err != nil {
				p.logger.Error("failed to update sink",
					zap.String("sinkID", sink.ID), zap.Error(err))
				continue
			}
		}
	}
	return
}
