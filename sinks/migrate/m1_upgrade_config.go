package migrate

import (
	"context"
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks"
	"github.com/orb-community/orb/sinks/authentication_type"
	"go.uber.org/zap"
)

type Plan1UpdateConfiguration struct {
	logger          *zap.Logger
	service         sinks.SinkService
	sinkRepo        sinks.SinkRepository
	passwordService authentication_type.PasswordService
}

func NewPlan1(logger *zap.Logger, service sinks.SinkService, sinkRepo sinks.SinkRepository, passwordService authentication_type.PasswordService) Plan {
	return &Plan1UpdateConfiguration{
		logger:          logger,
		service:         service,
		sinkRepo:        sinkRepo,
		passwordService: passwordService,
	}
}

func (p *Plan1UpdateConfiguration) Version() string {
	return "0.25.1"
}

func (p *Plan1UpdateConfiguration) Up(ctx context.Context) (mainErr error) {
	allSinks, mainErr := p.sinkRepo.SearchAllSinks(ctx, sinks.Filter{StateFilter: "", OpenTelemetry: "enabled"})
	if mainErr != nil {
		p.logger.Error("could not list sinks", zap.Error(mainErr))
		return
	}
	needsUpdate := 0
	updated := 0
	for _, sink := range allSinks {
		if _, ok := sink.Config["authentication"]; !ok {
			needsUpdate++
			sinkRemoteHost, ok := sink.Config["remote_host"]
			if !ok {
				p.logger.Error("failed to update sink for lack of remote_host", zap.String("sinkID", sink.ID))
				mainErr = errors.New("failed to get remote_host from config")
				continue
			}
			sinkUsername, ok := sink.Config["username"]
			if !ok {
				p.logger.Error("failed to update sink for lack of username", zap.String("sinkID", sink.ID))
				mainErr = errors.New("failed to get username from config")
				continue
			}
			encodedPassword, ok := sink.Config["password"]
			if !ok {
				p.logger.Error("failed to update sink for lack of password", zap.String("sinkID", sink.ID))
				mainErr = errors.New("failed to get password from config")
				continue
			}
			decodedPassword, err := p.passwordService.DecodePassword(encodedPassword.(string))
			if err != nil {
				p.logger.Error("failed to update sink for failure in decoding password",
					zap.String("sinkID", sink.ID), zap.Error(err))
				mainErr = err
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
				"opentelemetry": "enabled",
			}
			sink.Config = newMetadata
			_, err = p.service.UpdateSinkInternal(ctx, sink)
			if err != nil {
				p.logger.Error("failed to update sink",
					zap.String("sinkID", sink.ID), zap.Error(err))
				mainErr = err
				continue
			}
			updated++
		}
	}
	p.logger.Info("migration results", zap.Int("total_sinks", needsUpdate), zap.Int("updated_sinks", updated))
	return
}
