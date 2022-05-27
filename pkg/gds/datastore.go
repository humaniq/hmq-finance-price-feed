package gds

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
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

type Filter struct {
	Str   string
	Value interface{}
}

func ReadMultipleByFilters[T any](ctx context.Context, gds *Client, filters []Filter) ([]T, error) {
	query := datastore.NewQuery(gds.kind)
	for _, filter := range filters {
		query = query.Filter(filter.Str, filter.Value)
	}
	it := gds.ds.Run(ctx, query)
	var result []T
	for {
		var item T
		_, err := it.Next(&item)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrError, err)
		}
		result = append(result, item)
	}
	return result, nil
}
func WriteMultiple(ctx context.Context, gds *Client, values map[string]interface{}) error {
	keys := make([]*datastore.Key, 0, len(values))
	valuesList := make([]interface{}, 0, len(values))
	for key, val := range values {
		keys = append(keys, dsRecordKey(gds.kind, key))
		valuesList = append(valuesList, val)
	}
	if _, err := gds.ds.PutMulti(ctx, keys, valuesList); err != nil {
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
