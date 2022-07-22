package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type LRURecord struct {
	ExpiresAt time.Time
	Record    []byte
}

type LRU struct {
	cache  *lru.Cache
	expiry time.Duration
}

func NewLRU(size int) (*LRU, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &LRU{
		cache: cache,
	}, nil
}
func (l *LRU) Get(ctx context.Context, key string, value interface{}) error {
	val, found := l.cache.Get(key)
	if !found {
		return ErrNotFound
	}
	rec, ok := val.(*LRURecord)
	if !ok {
		return fmt.Errorf("%w: cast error", ErrReading)
	}
	if !rec.ExpiresAt.IsZero() && rec.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("%w: expired", ErrNotFound)
	}
	if err := json.Unmarshal(rec.Record, &value); err != nil {
		return fmt.Errorf("%w: %s", ErrDecoding, err)
	}
	return nil
}
func (l *LRU) Set(ctx context.Context, key string, value interface{}, expiryPeriod time.Duration) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(value); err != nil {
		return fmt.Errorf("%w: %s", ErrEncoding, err)
	}
	record := &LRURecord{Record: buf.Bytes()}
	if expiryPeriod != 0 {
		record.ExpiresAt = time.Now().Add(expiryPeriod)
	}
	if expiryPeriod == 0 && l.expiry != 0 {
		record.ExpiresAt = time.Now().Add(l.expiry)
	}
	_ = l.cache.Add(key, record)
	return nil
}
func (l *LRU) Unset(ctx context.Context, key string) error {
	l.cache.Remove(key)
	return nil
}
