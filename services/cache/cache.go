package cache

import (
	"errors"
	"sync"

	"github.com/mikestefanello/hooks"
	"github.com/mikestefanello/hooks-example/pkg/app"
	"github.com/samber/do"
)

type (
	// Cache provides a cache backend
	Cache interface {
		// Get gets cached data of a given key
		Get(key string) (any, error)

		// Set sets data in the cache with a given key
		Set(key string, data any) error

		// Delete deletes data from the cache with a given key
		Delete(key string) error
	}

	cache struct {
		store sync.Map
	}
)

func init() {
	// Provide dependencies during app boot process
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide(e.Msg, NewCache)
	})
}

// NewCache creates a new Cache instance
func NewCache(i *do.Injector) (Cache, error) {
	return &cache{}, nil
}

// Get gets cached data of a given key
func (c *cache) Get(key string) (any, error) {
	data, exists := c.store.Load(key)
	if !exists {
		return nil, errors.New("key does not exist")
	}
	return data, nil
}

// Set sets data in the cache with a given key
func (c *cache) Set(key string, data any) error {
	c.store.Store(key, data)
	return nil
}

// Delete deletes data from the cache with a given key
func (c *cache) Delete(key string) error {
	c.store.Delete(key)
	return nil
}
