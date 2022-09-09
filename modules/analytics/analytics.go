package analytics

import (
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/mikestefanello/hooks"
	"github.com/mikestefanello/hooks-example/modules/todo"
	"github.com/mikestefanello/hooks-example/pkg/app"
	"github.com/mikestefanello/hooks-example/services/web"
	"github.com/samber/do"
)

type (
	Analytics struct {
		WebRequests int64 `json:"webRequests"`
		Entities    int64 `json:"entities"`
	}

	Service interface {
		GetAnalytics() (Analytics, error)
		IncrementWebRequests() error
		IncrementEntities() error
	}

	Handler interface {
		Get(echo.Context) error
		WebRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	}

	analyticsService struct {
		analytics Analytics
		mu        sync.RWMutex
	}

	analyticsHandler struct {
		service Service
	}
)

var (
	HookAnalyticsUpdate = hooks.NewHook[Analytics]("analytics.update")
)

func init() {
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide(e.Msg, NewAnalyticsService)
		do.Provide(e.Msg, NewAnalyticsHandler)
	})

	web.HookBuildRouter.Listen(func(e hooks.Event[*echo.Echo]) {
		h := do.MustInvoke[Handler](do.DefaultInjector)
		e.Msg.GET("/analytics", h.Get)
		e.Msg.Use(h.WebRequestMiddleware)
	})

	todo.HookTodoInsert.Listen(func(e hooks.Event[todo.Todo]) {
		h := do.MustInvoke[Service](do.DefaultInjector)
		if err := h.IncrementEntities(); err != nil {
			log.Error(err)
		}
	})
}
