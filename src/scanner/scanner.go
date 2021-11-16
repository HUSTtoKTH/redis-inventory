// Package scanner TODO
package scanner

import (
	"context"

	"github.com/obukhov/redis-inventory/src/adapter"
	"github.com/obukhov/redis-inventory/src/trie"
	"github.com/obukhov/redis-inventory/src/typetrie.go"
	"github.com/rs/zerolog"
)

// RedisServiceInterface abstraction to access redis
type RedisServiceInterface interface {
	ScanKeys(ctx context.Context, options adapter.ScanOptions) <-chan *adapter.KeyInfo
	GetKeysCount(ctx context.Context) (int64, error)
	GetMemoryUsage(ctx context.Context, key adapter.KeyInfo) (int64, error)
	GetKeyType(ctx context.Context, key *adapter.KeyInfo)
	GetTypeBatch(ctx context.Context, keys []*adapter.KeyInfo)
	GetMemoryUsageBatch(ctx context.Context, keys []*adapter.KeyInfo)
}

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
func (s *RedisScanner) Scan(options adapter.ScanOptions, result *typetrie.TypeTrie) {
	var totalCount int64
	if options.Pattern == "*" || options.Pattern == "" {
		totalCount = s.getKeysCount()
	}

	s.scanProgress.Start(totalCount)
	keys := []*adapter.KeyInfo{}
	for key := range s.redisService.ScanKeys(context.Background(), options) {
		s.scanProgress.Increment()
		keys = append(keys, key)
		if len(keys) == 1000 {
			s.redisService.GetMemoryUsageBatch(context.Background(), keys)
			s.redisService.GetTypeBatch(context.Background(), keys)
			for _, key := range keys {
				result.Add(
					key.Key,
					key.Type,
					trie.ParamValue{Param: trie.BytesSize, Value: key.BytesSize},
					trie.ParamValue{Param: trie.KeysCount, Value: 1},
				)
			}
			keys = []*adapter.KeyInfo{}
		}
	}
	s.redisService.GetMemoryUsageBatch(context.Background(), keys)
	s.redisService.GetTypeBatch(context.Background(), keys)
	for _, key := range keys {
		result.Add(
			key.Key,
			key.Type,
			trie.ParamValue{Param: trie.BytesSize, Value: key.BytesSize},
			trie.ParamValue{Param: trie.KeysCount, Value: 1},
		)
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
