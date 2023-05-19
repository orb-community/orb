package migrate

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/orb-community/orb/sinks"
	"go.uber.org/zap"
)

type Plan interface {
	Version() string
	Up(ctx context.Context) error
}

type Service interface {
	Migrate(plans ...Plan) error
}

func NewService(logger *zap.Logger, sinkRepository sinks.SinkRepository) Service {
	return &migrateService{
		logger:         logger,
		sinkRepository: sinkRepository,
	}
}

type migrateService struct {
	logger         *zap.Logger
	sinkRepository sinks.SinkRepository
}

func (m *migrateService) updateNewVersion(ctx context.Context, newVersion string) {
	currentVersion, err := m.getCurrentVersion(ctx)
	if err != nil {
		m.logger.Error("Could not parse current version", zap.Error(err))
		return
	}
	incomingSemVer, err := version.NewSemver(newVersion)
	if err != nil {
		m.logger.Error("Could not parse version for plan", zap.String("version", newVersion), zap.Error(err))
		return
	}
	if incomingSemVer.GreaterThan(currentVersion) {
		err := m.sinkRepository.UpdateVersion(ctx, newVersion)
		if err != nil {
			m.logger.Error("error during update of version", zap.String("newVersion", newVersion), zap.Error(err))
			return
		}
	}
}

func (m *migrateService) getCurrentVersion(ctx context.Context) (*version.Version, error) {
	currentVersion, _ := m.sinkRepository.GetVersion(ctx)
	return version.NewSemver(currentVersion)
}

func (m *migrateService) Migrate(plans ...Plan) error {
	for i, plan := range plans {
		planName := fmt.Sprintf("plan%d", i)
		ctx := context.WithValue(context.Background(), "migrate", planName)
		v, err := version.NewSemver(plan.Version())
		if err != nil {
			m.logger.Error("Could not parse version for plan", zap.String("version", plan.Version()), zap.Error(err))
		}
		curV, err := m.getCurrentVersion(ctx)
		if err != nil {
			m.logger.Error("could not find current version, version", zap.Error(err))
			return err
		}
		if v.GreaterThan(curV) {
			m.logger.Info("Starting plan", zap.Int("plan", i))
			err := plan.Up(ctx)
			if err != nil {
				m.logger.Error("error during migrate service", zap.Error(err))
				return err
			}
			m.updateNewVersion(ctx, plan.Version())
		}
	}
	return nil
}
