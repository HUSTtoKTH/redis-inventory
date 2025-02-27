package adapter

import (
	"context"
	"time"

	"github.com/mediocregopher/radix/v4"
)

// KeyInfo TODO
type KeyInfo struct {
	Key       string
	Node      string
	Type      string
	BytesSize int64
}

// ScanOptions options for scanning keyspace
type ScanOptions struct {
	Pattern   string
	ScanCount int
	Throttle  int
}

// NewRedisService creates RedisService
func NewRedisService(client radix.Client) RedisService {
	return RedisService{
		client: client,
	}
}

// RedisService implementation for iteration over redis
type RedisService struct {
	client radix.Client
}

// ScanKeys scans keys asynchroniously and sends them to the returned channel
func (s RedisService) ScanKeys(ctx context.Context, options ScanOptions) <-chan *KeyInfo {
	resultChan := make(chan *KeyInfo)

	scanOpts := radix.ScannerConfig{
		Command: "SCAN",
		Count:   options.ScanCount,
	}

	if options.Pattern != "*" && options.Pattern != "" {
		scanOpts.Pattern = options.Pattern
	}

	go func() {
		defer close(resultChan)
		var key string
		radixScanner := scanOpts.New(s.client)
		for radixScanner.Next(ctx, &key) {
			resultChan <- &KeyInfo{
				Key: key,
			}
			if options.Throttle > 0 {
				time.Sleep(time.Nanosecond * time.Duration(options.Throttle))
			}
		}
	}()

	return resultChan
}

// GetKeysCount returns number of keys in the current database
func (s RedisService) GetKeysCount(ctx context.Context) (int64, error) {
	var keysCount int64
	err := s.client.Do(context.Background(), radix.Cmd(&keysCount, "DBSIZE"))
	if err != nil {
		return 0, err
	}

	return keysCount, nil
}

// GetMemoryUsage returns memory usage of given key
func (s RedisService) GetMemoryUsage(ctx context.Context, key KeyInfo) (int64, error) {
	var res int64
	err := s.client.Do(context.Background(), radix.Cmd(&res, "MEMORY", "USAGE", key.Key))
	if err != nil {
		return 0, err
	}

	return res, nil
}

// GetKeyType TODO
func (s RedisService) GetKeyType(ctx context.Context, key *KeyInfo) {
	var res string
	err := s.client.Do(context.Background(), radix.Cmd(&res, "TYPE", "USAGE", key.Key))
	if err != nil {
		return
	}
	key.Type = res
	return
}
