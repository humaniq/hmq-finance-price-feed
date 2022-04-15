package gds

import (
	"context"

	"cloud.google.com/go/datastore"
)

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
