package resource

import (
	"context"
	"testing"
	"time"

	v1sync "github.com/apache/servicecomb-service-center/syncer/api/v1"

	"github.com/go-chassis/cari/sync"
	"github.com/stretchr/testify/assert"
)

func TestNewKV(t *testing.T) {
	// create kv case
	id, _ := v1sync.NewEventID()
	e := &v1sync.Event{
		Id:      id,
		Action:  sync.CreateAction,
		Subject: KV,
		Opts: map[string]string{
			"key": "hello",
		},
		Value:     []byte("world"),
		Timestamp: v1sync.Timestamp(),
	}
	fm := &mockKVManager{
		kvs: make(map[string][]byte, 10),
	}
	k := &kv{
		event: e,

		manager: fm,
	}
	ctx := context.Background()
	result := k.LoadCurrentResource(ctx)
	if assert.Nil(t, result) {
		result = k.NeedOperate(ctx)
		if assert.Nil(t, result) {
			result = k.Operate(ctx)
			if assert.NotNil(t, result) && assert.Equal(t, Success, result.Status) {
				data, err := k.manager.Get(ctx, "hello")
				assert.Nil(t, err)
				assert.Equal(t, []byte("world"), data)
			}
		}
	}

	e1 := &v1sync.Event{
		Id:      id,
		Action:  sync.UpdateAction,
		Subject: KV,
		Opts: map[string]string{
			"key": "hello1",
		},
		Value:     []byte("change"),
		Timestamp: v1sync.Timestamp(),
	}
	// Tombstone is later than event case
	k = &kv{
		event: e1,

		manager: fm,
		tombstoneLoader: &mockTombstoneLoader{
			ts: &sync.Tombstone{
				ResourceID:   "xxx2",
				ResourceType: "kv",
				Timestamp:    time.Now().Add(time.Minute).UnixNano(),
			},
		},
	}
	result = k.LoadCurrentResource(ctx)
	if assert.Nil(t, result) {
		result = k.NeedOperate(ctx)
		if assert.NotNil(t, result) {
			assert.Equal(t, Skip, result.Status)
		}
	}

	// Tombstone is not exist case
	k = &kv{
		event: e1,

		manager: fm,
		tombstoneLoader: &mockTombstoneLoader{
			ts: nil,
		},
	}

	result = k.LoadCurrentResource(ctx)
	if assert.Nil(t, result) {
		result = k.NeedOperate(ctx)
		assert.Nil(t, result)
		result = k.Operate(ctx)
		if assert.NotNil(t, result) && assert.Equal(t, Success, result.Status) {
			data, err := k.manager.Get(ctx, "hello1")
			assert.Nil(t, err)
			assert.Equal(t, []byte("change"), data)
		}
	}

	// delete case
	e2 := &v1sync.Event{
		Id:      id,
		Action:  sync.DeleteAction,
		Subject: KV,
		Opts: map[string]string{
			"key": "hello",
		},
		Value:     []byte("change"),
		Timestamp: v1sync.Timestamp(),
	}
	k = &kv{
		event: e2,

		manager: fm,
		tombstoneLoader: &mockTombstoneLoader{
			ts: nil,
		},
	}
	result = k.LoadCurrentResource(ctx)
	if assert.Nil(t, result) {
		result = k.NeedOperate(ctx)
		assert.Nil(t, result)
	}
	result = k.Operate(ctx)
	if assert.NotNil(t, result) && assert.Equal(t, Success, result.Status) {
		data, err := k.manager.Get(ctx, "hello")
		assert.Equal(t, ErrRecordNonExist, err)
		assert.Nil(t, data)
	}
}

type mockKVManager struct {
	kvs map[string][]byte
}

func (f *mockKVManager) Get(_ context.Context, key string) ([]byte, error) {
	result, ok := f.kvs[key]
	if !ok {
		return nil, ErrRecordNonExist
	}
	return result, nil
}

func (f *mockKVManager) Put(_ context.Context, key string, value []byte) error {
	f.kvs[key] = value
	return nil
}

func (f *mockKVManager) Post(_ context.Context, key string, value []byte) error {
	f.kvs[key] = value
	return nil
}

func (f *mockKVManager) Delete(_ context.Context, key string) error {
	_, ok := f.kvs[key]
	if !ok {
		return nil
	}
	delete(f.kvs, key)
	return nil
}
