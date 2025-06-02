package service_test

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"protokey/pkg/app/service"
	"protokey/pkg/infrastructure/storage"
)

func createTempFile(t *testing.T) string {
	file, err := os.CreateTemp("", "store_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	path := file.Name()
	_ = file.Close()

	return path
}

func cleanupFile(path string) {
	_ = os.Remove(path)
}

func TestBasicSetGet(t *testing.T) {
	path := createTempFile(t)
	defer cleanupFile(path)

	store := storage.NewStore(storage.Config{DataFile: path})
	protoKeyService := service.NewProtoKeyService(store.CommandChan)

	value, err := protoKeyService.Get("nonexistent")
	assert.NoError(t, err)
	assert.Equal(t, 0, value)

	err = protoKeyService.Set("existent", 42)
	assert.NoError(t, err)

	value, err = protoKeyService.Get("existent")
	assert.NoError(t, err)
	assert.Equal(t, 42, value)

	err = protoKeyService.Set("existent", 0)
	assert.NoError(t, err)

	value, err = protoKeyService.Get("existent")
	assert.NoError(t, err)
	assert.Equal(t, 0, value)

	close(store.CommandChan)
}

func TestKeysAndList(t *testing.T) {
	path := createTempFile(t)
	defer cleanupFile(path)

	store := storage.NewStore(storage.Config{DataFile: path})
	protoKeyService := service.NewProtoKeyService(store.CommandChan)

	assert.NoError(t, protoKeyService.Set("a1", 1))
	assert.NoError(t, protoKeyService.Set("a2", 2))
	assert.NoError(t, protoKeyService.Set("b1", 3))
	assert.NoError(t, protoKeyService.Set("a10", 10))

	keys, err := protoKeyService.Keys("a")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"a1", "a2", "a10"}, keys)

	keys, err = protoKeyService.Keys("c")
	assert.NoError(t, err)
	assert.Empty(t, keys)

	close(store.CommandChan)
}

func TestInvalidKey(t *testing.T) {
	path := createTempFile(t)
	defer cleanupFile(path)

	store := storage.NewStore(storage.Config{DataFile: path})
	protoKeyService := service.NewProtoKeyService(store.CommandChan)

	err := protoKeyService.Set("bad key", 5)
	assert.ErrorIs(t, err, service.ErrBadRequest)

	err = protoKeyService.Set(strings.Repeat("a", 10001), 5)
	assert.ErrorIs(t, err, service.ErrBadRequest)

	_, err = protoKeyService.Get("?!")
	assert.ErrorIs(t, err, service.ErrBadRequest)

	_, err = protoKeyService.Keys("   ")
	assert.ErrorIs(t, err, service.ErrBadRequest)

	close(store.CommandChan)
}

func TestConcurrentAccess(t *testing.T) {
	path := createTempFile(t)
	defer cleanupFile(path)

	store := storage.NewStore(storage.Config{DataFile: path})
	protoKeyService := service.NewProtoKeyService(store.CommandChan)

	const N = 100
	var waitGroup sync.WaitGroup
	waitGroup.Add(N)

	for i := 0; i < N; i++ {
		go func(i int) {
			defer waitGroup.Done()
			key := "key_" + strconv.Itoa(i)
			assert.NoError(t, protoKeyService.Set(key, i*10))
			value, err := protoKeyService.Get(key)
			assert.NoError(t, err)
			assert.Equal(t, i*10, value)
		}(i)
	}

	waitGroup.Wait()

	for i := 0; i < N; i++ {
		key := "key_" + strconv.Itoa(i)
		value, err := protoKeyService.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, i*10, value)
	}

	close(store.CommandChan)
}

func TestPersistence(t *testing.T) {
	path := createTempFile(t)
	defer cleanupFile(path)

	store1 := storage.NewStore(storage.Config{DataFile: path})
	protoKeyService1 := service.NewProtoKeyService(store1.CommandChan)

	assert.NoError(t, protoKeyService1.Set("first", 11))
	assert.NoError(t, protoKeyService1.Set("second", 22))

	time.Sleep(1200 * time.Millisecond)

	close(store1.CommandChan)
	time.Sleep(100 * time.Millisecond)

	store2 := storage.NewStore(storage.Config{DataFile: path})
	protoKeyService2 := service.NewProtoKeyService(store2.CommandChan)

	value1, err := protoKeyService2.Get("first")
	assert.NoError(t, err)
	assert.Equal(t, 11, value1)

	value2, err := protoKeyService2.Get("second")
	assert.NoError(t, err)
	assert.Equal(t, 22, value2)

	close(store2.CommandChan)
}
