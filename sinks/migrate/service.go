package migrate

import (
	"context"
	"fmt"
	"github.com/orb-community/orb/pkg/types"
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
	currentVersion := m.getCurrentVersion(ctx)
	incomingSemVer := types.NewSemVerFromString(newVersion)
	if !currentVersion.IsEqualTo(incomingSemVer) && incomingSemVer.IsNewerThan(currentVersion) {
		err := m.sinkRepository.UpdateVersion(ctx, newVersion)
		if err != nil {
			m.logger.Error("error during update of version", zap.String("newVersion", newVersion), zap.Error(err))
			return
		}
	}
}

func (m *migrateService) getCurrentVersion(ctx context.Context) types.SemVer {
	currentVersion, _ := m.sinkRepository.GetVersion(ctx)
	currSemVer := types.NewSemVerFromString(currentVersion)
	return currSemVer
}

func (m *migrateService) Migrate(plans ...Plan) error {

	for i, plan := range plans {
		planName := fmt.Sprintf("plan%d", i)
		ctx := context.WithValue(context.Background(), "migrate", planName)

		m.logger.Info("Starting plan", zap.Int("plan", i))
		err := plan.Up(ctx)
		if err != nil {
			m.logger.Error("error during migrate service", zap.Error(err))
			return err
		}
		m.updateNewVersion(ctx, plan.Version())
	}
	return nil
}
