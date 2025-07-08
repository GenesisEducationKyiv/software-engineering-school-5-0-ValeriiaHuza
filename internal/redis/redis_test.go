//go:build unit
// +build unit

package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock redisClient ---

type mockRedisClient struct {
	mock.Mock
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, ttl)
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetErr(args.Error(0))
	return cmd
}

func (m *mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx)
	if err := args.Error(1); err != nil {
		cmd.SetErr(err)
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

func (m *mockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	cmd := redis.NewIntCmd(ctx)
	cmd.SetErr(args.Error(0))
	return cmd
}

func (m *mockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	return cmd
}

// --- Tests ---

func TestSet_Success(t *testing.T) {
	mockClient := new(mockRedisClient)
	ctx := context.Background()
	provider := NewRedisProvider(mockClient, ctx)

	val := map[string]string{"foo": "bar"}
	data, _ := json.Marshal(val)
	mockClient.On("Set", ctx, "key1", data, time.Duration(0)).Return(nil)

	err := provider.Set("key1", val)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSetWithTTL_Success(t *testing.T) {
	mockClient := new(mockRedisClient)
	ctx := context.Background()
	provider := NewRedisProvider(mockClient, ctx)

	val := map[string]string{"foo": "bar"}
	data, _ := json.Marshal(val)
	ttl := 5 * time.Minute
	mockClient.On("Set", ctx, "key2", data, ttl).Return(nil)

	err := provider.SetWithTTL("key2", val, ttl)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestGet_Success(t *testing.T) {
	mockClient := new(mockRedisClient)
	ctx := context.Background()
	provider := NewRedisProvider(mockClient, ctx)

	val := map[string]string{"foo": "bar"}
	data, _ := json.Marshal(val)
	mockClient.On("Get", ctx, "key3").Return(string(data), nil)

	var result map[string]string
	err := provider.Get("key3", &result)
	assert.NoError(t, err)
	assert.Equal(t, val, result)
	mockClient.AssertExpectations(t)
}

func TestGet_Error(t *testing.T) {
	mockClient := new(mockRedisClient)
	ctx := context.Background()
	provider := NewRedisProvider(mockClient, ctx)

	mockClient.On("Get", ctx, "key4").Return("", errors.New("not found"))

	var result map[string]string
	err := provider.Get("key4", &result)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockClient := new(mockRedisClient)
	ctx := context.Background()
	provider := NewRedisProvider(mockClient, ctx)

	mockClient.On("Del", ctx, []string{"key5"}).Return(nil)

	err := provider.Delete("key5")
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestDelete_Error(t *testing.T) {
	mockClient := new(mockRedisClient)
	ctx := context.Background()
	provider := NewRedisProvider(mockClient, ctx)

	mockClient.On("Del", ctx, []string{"key6"}).Return(errors.New("del error"))

	err := provider.Delete("key6")
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}
