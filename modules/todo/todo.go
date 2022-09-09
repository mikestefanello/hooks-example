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
	// Todo is the todo model
	Todo struct {
		ID       int    `json:"id"`
		Label    string `json:"label"`
		Complete bool   `json:"complete"`
	}

	// Service provides an interface to service todos
	Service interface {
		// GetTodo loads a todo by ID
		GetTodo(id int) (Todo, error)

		// GetTodos loads all todos
		GetTodos() ([]Todo, error)

		// InsertTodo creates and saves a new todo
		InsertTodo(*Todo) error
	}

	// Handler provides an HTTP handler for todos
	Handler interface {
		// Index handles a request for all todos
		Index(echo.Context) error

		// Get handles a request for a single todo
		Get(echo.Context) error

		// Post handles a request to create a new todo
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
	// HookTodoPreInsert allows modules the ability to alter todos prior to saving them
	HookTodoPreInsert = hooks.NewHook[*Todo]("todo.pre_insert")

	// HookTodoInsert allows modules to listen for when todos are inserted
	HookTodoInsert = hooks.NewHook[Todo]("todo.insert")
)

func init() {
	// Provide dependencies during app boot process
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide(e.Msg, NewTodoService)
		do.Provide(e.Msg, NewTodoHandler)
	})

	// Provide web routes
	web.HookBuildRouter.Listen(func(e hooks.Event[*echo.Echo]) {
		h := do.MustInvoke[Handler](do.DefaultInjector)
		e.Msg.GET("/todo", h.Index)
		e.Msg.GET("/todo/:todo", h.Get)
		e.Msg.POST("/todo", h.Post)
	})
}
