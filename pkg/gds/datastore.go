package gds

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/datastore"
)

var ErrNotFound = errors.New("not found")
var ErrError = errors.New("error")

type Client struct {
	ds   *datastore.Client
	kind string
}

func NewClient(ctx context.Context, projectKey string, kind string) (*Client, error) {
	client, err := datastore.NewClient(ctx, projectKey)
	if err != nil {
		return nil, err
	}
	return &Client{ds: client, kind: kind}, nil
}
func (ds *Client) Write(ctx context.Context, key string, record interface{}) error {
	if err := dsWriteRecord(ctx, ds.ds, ds.kind, key, record); err != nil {
		return fmt.Errorf("%w: %s", ErrError, err)
	}
	return nil
}
func (ds *Client) Read(ctx context.Context, key string, record interface{}) error {
	if err := dsReadRecord(ctx, ds.ds, ds.kind, key, record); err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return fmt.Errorf("%w: %s", ErrNotFound, err)
		}
		return fmt.Errorf("%w: %s", ErrError, err)
	}
	return nil
}

func dsWriteRecord(ctx context.Context, ds *datastore.Client, kind string, key string, record interface{}) error {
	_, err := ds.Put(ctx, dsRecordKey(kind, key), record)
	if err != nil {
		return err
	}
	return nil
}
func dsReadRecord(ctx context.Context, ds *datastore.Client, kind string, key string, record interface{}) error {
	if err := ds.Get(ctx, dsRecordKey(kind, key), record); err != nil {
		return err
	}
	return nil
}

func dsRecordKey(kind string, key string) *datastore.Key {
	return datastore.NameKey(kind, key, nil)
}
