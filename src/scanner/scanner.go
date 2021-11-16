// Package scanner TODO
package scanner

import (
	"context"

	"github.com/obukhov/redis-inventory/src/adapter"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/rs/zerolog"
)

// RedisServiceInterface abstraction to access redis
type RedisServiceInterface interface {
	ScanKeys(ctx context.Context, options adapter.ScanOptions) <-chan adapter.KeyInfo
	GetKeysCount(ctx context.Context) (int64, error)
	GetMemoryUsage(ctx context.Context, key adapter.KeyInfo) (int64, error)
}

// KeyInfo TODO

// RedisScanner scans redis keys and puts them in a trie
type RedisScanner struct {
	redisService RedisServiceInterface
	scanProgress adapter.ProgressWriter
	logger       zerolog.Logger
}

// NewScanner creates RedisScanner
func NewScanner(redisService RedisServiceInterface, scanProgress adapter.ProgressWriter, logger zerolog.Logger) *RedisScanner {
	return &RedisScanner{
		redisService: redisService,
		scanProgress: scanProgress,
		logger:       logger,
	}
}

// Scan initiates scanning process
func (s *RedisScanner) Scan(options adapter.ScanOptions, result *trie.Trie) {
	var totalCount int64
	if options.Pattern == "*" || options.Pattern == "" {
		totalCount = s.getKeysCount()
	}

	s.scanProgress.Start(totalCount)
	for key := range s.redisService.ScanKeys(context.Background(), options) {
		s.scanProgress.Increment()
		res, err := s.redisService.GetMemoryUsage(context.Background(), key)
		if err != nil {
			s.logger.Error().Err(err).Msgf("Error dumping key %s", key)
			continue
		}

		result.Add(
			key.Key,
			trie.ParamValue{Param: trie.BytesSize, Value: res},
			trie.ParamValue{Param: trie.KeysCount, Value: 1},
		)

		s.logger.Debug().Msgf("Dump %s value: %d", key, res)
	}
	s.scanProgress.Stop()
}

func (s *RedisScanner) getKeysCount() int64 {
	res, err := s.redisService.GetKeysCount(context.Background())
	if err != nil {
		s.logger.Error().Err(err).Msgf("Error getting number of keys")
		return 0
	}
	s.logger.Info().Msgf("key number: %d", res)
	return res
}
