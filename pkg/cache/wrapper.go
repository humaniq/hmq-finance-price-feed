package cache

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")
var ErrEncoding = errors.New("error encoding")
var ErrDecoding = errors.New("error decoding")
var ErrReading = errors.New("error reading")
var ErrWriting = errors.New("error writing")

type Wrapper interface {
	Set(ctx context.Context, key string, value interface{}, expiryPeriod time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
}
