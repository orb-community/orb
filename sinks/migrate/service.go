package migrate

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

type Plan interface {
	Version() string
	Up(ctx context.Context) error
}

type Service interface {
	Migrate(plans ...Plan) error
}

func NewService(logger *zap.Logger) Service {
	return &migrateService{
		logger: logger,
	}
}

type migrateService struct {
	logger *zap.Logger
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
	}

}
