package producer

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type SinkerKey struct {
	OwnerID      string
	SinkID       string
	Size         string
	LastActivity time.Time
}

func (s *SinkerKey) Encode() map[string]interface{} {
	return map[string]interface{}{
		"owner_id":      s.OwnerID,
		"sink_id":       s.SinkID,
		"size":          s.Size,
		"last_activity": s.LastActivity.Format(time.RFC3339),
	}
}

const DefaultExpiration = 5 * time.Minute

type SinkerKeyService interface {
	// AddNewSinkerKey Add New Sinker Key with default Expiration of 5 minutes
	AddNewSinkerKey(ctx context.Context, key SinkerKey) error
	// RenewSinkerKey Increment Expiration of Sinker Key
	RenewSinkerKey(ctx context.Context, key SinkerKey) error
	// RenewSinkerKeyInternal Increment Expiration of Sinker Key
	RenewSinkerKeyInternal(ctx context.Context, sink SinkerKey, expiration time.Duration) error
}

type sinkerKeyService struct {
	logger          *zap.Logger
	cacheRepository *redis.Client
}

func NewSinkerKeyService(l *zap.Logger, cacheRepository *redis.Client) SinkerKeyService {
	logger := l.Named("sinker_key_service")
	return &sinkerKeyService{logger: logger, cacheRepository: cacheRepository}
}

// RenewSinkerKey Increment Expiration of Sinker Key
func (s *sinkerKeyService) RenewSinkerKey(ctx context.Context, sink SinkerKey) error {
	// If key does not exist, create new entry
	key := fmt.Sprintf("orb.sinker.key-%s:%s", sink.OwnerID, sink.SinkID)
	cmd := s.cacheRepository.Expire(ctx, key, DefaultExpiration)
	if cmd.Err() != nil {
		s.logger.Error("error sending event to sinker event store", zap.Error(cmd.Err()))
		return cmd.Err()
	}
	return nil
}

// RenewSinkerKeyInternal Increment Expiration of Sinker Key using custom expiration
func (s *sinkerKeyService) RenewSinkerKeyInternal(ctx context.Context, sink SinkerKey, expiration time.Duration) error {
	// If key does not exist, create new entry
	key := fmt.Sprintf("orb.sinker.key-%s:%s", sink.OwnerID, sink.SinkID)
	cmd := s.cacheRepository.Expire(ctx, key, expiration)
	if cmd.Err() != nil {
		s.logger.Error("error sending event to sinker event store", zap.Error(cmd.Err()))
		return cmd.Err()
	}
	return nil
}

func (s *sinkerKeyService) AddNewSinkerKey(ctx context.Context, sink SinkerKey) error {
	// Create sinker key in redis Hashset with default expiration of 5 minutes
	key := fmt.Sprintf("orb.sinker.key-%s:%s", sink.OwnerID, sink.SinkID)
	cmd := s.cacheRepository.HSet(ctx, key, sink.Encode())
	if cmd.Err() != nil {
		s.logger.Error("error sending event to sinker event store", zap.Error(cmd.Err()))
		return cmd.Err()
	}
	err := s.RenewSinkerKey(ctx, sink)
	if err != nil {
		s.logger.Error("error setting expiration to sinker event store", zap.Error(cmd.Err()))
		return cmd.Err()
	}
	return nil
}
