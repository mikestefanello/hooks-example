package analytics

import (
	"sync"

	"github.com/samber/do"
)

func NewAnalyticsService(i *do.Injector) (Service, error) {
	return &analyticsService{
		analytics: Analytics{},
		mu:        sync.RWMutex{},
	}, nil
}

func (a *analyticsService) GetAnalytics() (Analytics, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.analytics, nil
}

func (a *analyticsService) IncrementWebRequests() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.analytics.WebRequests++
	return nil
}

func (a *analyticsService) IncrementEntities() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.analytics.Entities++
	return nil
}
