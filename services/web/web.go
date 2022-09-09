package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mikestefanello/hooks"
	"github.com/mikestefanello/hooks-example/pkg/app"
	"github.com/mikestefanello/hooks-example/services/config"
	"github.com/samber/do"
)

type (
	// Web provides a web server
	Web interface {
		// Start starts the web server and will block until it is stopped
		Start() error
	}

	web struct {
		handler *echo.Echo
		cfg     config.Config
	}
)

// HookBuildRouter allows modules the ability to build on the web router
var HookBuildRouter = hooks.NewHook[*echo.Echo]("router.build")

func init() {
	// Provide dependencies during app boot process
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide(e.Msg, NewWeb)
	})
}

// NewWeb creates a new Web instance
func NewWeb(i *do.Injector) (Web, error) {
	w := &web{
		handler: echo.New(),
		cfg:     do.MustInvoke[config.Config](i),
	}
	w.buildRouter()

	return w, nil
}

// buildRouter builds the web router
func (w *web) buildRouter() {
	w.handler.Use(
		middleware.RequestID(),
		middleware.Logger(),
	)

	w.handler.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "hello world")
	})

	// Allow all modules to build on the router
	HookBuildRouter.Dispatch(w.handler)

	// Output the routes to the log
	routes := make([]string, len(w.handler.Routes()))
	for i, r := range w.handler.Routes() {
		routes[i] = fmt.Sprintf("%s_%s", r.Method, r.Path)
	}
	log.Printf("registered %d routes: %v", len(routes), routes)
}

// Start starts the web server and will block until it is stopped
func (w *web) Start() error {
	httpCfg := w.cfg.GetHTTP()

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", httpCfg.Hostname, httpCfg.Port),
		Handler:      w.handler,
		ReadTimeout:  httpCfg.ReadTimeout,
		WriteTimeout: httpCfg.WriteTimeout,
		IdleTimeout:  httpCfg.IdleTimeout,
	}

	return w.handler.StartServer(&srv)
}
