package todo

import (
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/hooks"
	"github.com/mikestefanello/hooks-example/pkg/app"
	"github.com/mikestefanello/hooks-example/services/cache"
	"github.com/mikestefanello/hooks-example/services/web"
	"github.com/samber/do"
)

type (
	Todo struct {
		ID       int    `json:"id"`
		Label    string `json:"label"`
		Complete bool   `json:"complete"`
	}

	Service interface {
		GetTodo(id int) (Todo, error)
		GetTodos() ([]Todo, error)
		InsertTodo(*Todo) error
	}

	Handler interface {
		Index(echo.Context) error
		Get(echo.Context) error
		Post(echo.Context) error
	}

	todoService struct {
		cache cache.Cache
	}

	todoHandler struct {
		service Service
	}
)

var (
	HookTodoPreInsert = hooks.NewHook[*Todo]("todo.pre_insert")
	HookTodoInsert    = hooks.NewHook[Todo]("todo.insert")
)

func init() {
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide(e.Msg, NewTodoService)
		do.Provide(e.Msg, NewTodoHandler)
	})

	web.HookBuildRouter.Listen(func(e hooks.Event[*echo.Echo]) {
		h := do.MustInvoke[Handler](do.DefaultInjector)
		e.Msg.GET("/todo", h.Index)
		e.Msg.GET("/todo/:todo", h.Get)
		e.Msg.POST("/todo", h.Post)
	})
}
