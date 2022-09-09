package cache

import (
	"errors"
	"sync"

	"github.com/mikestefanello/hooks"
	"github.com/mikestefanello/hooks-example/pkg/app"
	"github.com/samber/do"
)

type (
	Cache interface {
		Get(key string) (any, error)
		Set(key string, data any) error
		Delete(key string) error
	}

	cache struct {
		store sync.Map
	}
)

func init() {
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide(e.Msg, NewCache)
	})
}

func NewCache(i *do.Injector) (Cache, error) {
	return &cache{}, nil
}

func (c *cache) Get(key string) (any, error) {
	data, exists := c.store.Load(key)
	if !exists {
		return nil, errors.New("key does not exist")
	}
	return data, nil
}

func (c *cache) Set(key string, data any) error {
	c.store.Store(key, data)
	return nil
}

func (c *cache) Delete(key string) error {
	c.store.Delete(key)
	return nil
}
