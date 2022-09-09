package analytics

import (
	"github.com/mikestefanello/hooks-example/services/cache"
	"github.com/samber/do"
)

const cacheKey = "analytics"

func NewAnalyticsService(i *do.Injector) (Service, error) {
	return &analyticsService{
		cache: do.MustInvoke[cache.Cache](i),
	}, nil
}

func (a *analyticsService) GetAnalytics() (Analytics, error) {
	return a.load(), nil
}

func (a *analyticsService) IncrementWebRequests() error {
	data := a.load()
	data.WebRequests++

	if err := a.save(data); err != nil {
		return nil
	}

	HookAnalyticsUpdate.Dispatch(data)

	return nil
}

func (a *analyticsService) IncrementEntities() error {
	data := a.load()
	data.Entities++

	if err := a.save(data); err != nil {
		return nil
	}

	HookAnalyticsUpdate.Dispatch(data)

	return nil
}

func (a *analyticsService) load() Analytics {
	data, err := a.cache.Get(cacheKey)
	if err != nil {
		return Analytics{}
	}

	return data.(Analytics)
}

func (a *analyticsService) save(data Analytics) error {
	return a.cache.Set(cacheKey, data)
}
