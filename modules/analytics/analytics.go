package analytics

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/mikestefanello/hooks"
	"github.com/mikestefanello/hooks-example/modules/todo"
	"github.com/mikestefanello/hooks-example/pkg/app"
	"github.com/mikestefanello/hooks-example/services/cache"
	"github.com/mikestefanello/hooks-example/services/web"
	"github.com/samber/do"
)

type (
	// Analytics is the analytics model
	Analytics struct {
		WebRequests int64 `json:"webRequests"`
		Entities    int64 `json:"entities"`
	}

	// Service provides an interface to service analytics
	Service interface {
		// GetAnalytics returns analytics
		GetAnalytics() (Analytics, error)

		// IncrementWebRequests increments the counter that tracks web requests
		IncrementWebRequests() error

		// IncrementEntities increments the counter that track entities
		IncrementEntities() error
	}

	// Handler provides an HTTP handler for analytics
	Handler interface {
		// Get handles web requests to get analytics
		Get(echo.Context) error

		// WebRequestMiddleware provides middleware to track all web requests
		WebRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	}

	analyticsService struct {
		cache cache.Cache
	}

	analyticsHandler struct {
		service Service
	}
)

var (
	// HookAnalyticsUpdate allows modules to listen for analytics updates
	HookAnalyticsUpdate = hooks.NewHook[Analytics]("analytics.update")
)

func init() {
	// Provide dependencies during app boot process
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide(e.Msg, NewAnalyticsService)
		do.Provide(e.Msg, NewAnalyticsHandler)
	})

	// Provide web routes
	web.HookBuildRouter.Listen(func(e hooks.Event[*echo.Echo]) {
		h := do.MustInvoke[Handler](do.DefaultInjector)
		e.Msg.GET("/analytics", h.Get)
		e.Msg.Use(h.WebRequestMiddleware)
	})

	// React to new todos being inserted
	todo.HookTodoInsert.Listen(func(e hooks.Event[todo.Todo]) {
		h := do.MustInvoke[Service](do.DefaultInjector)
		if err := h.IncrementEntities(); err != nil {
			log.Error(err)
		}
	})
}
