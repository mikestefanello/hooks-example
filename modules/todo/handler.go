package todo

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

// NewTodoHandler provides a new todo Handler instance
func NewTodoHandler(i *do.Injector) (Handler, error) {
	return &todoHandler{
		service: do.MustInvoke[Service](i),
	}, nil
}

func (t *todoHandler) Index(ctx echo.Context) error {
	todos, err := t.service.GetTodos()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, todos)
}

func (t *todoHandler) Get(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("todo"))
	if err != nil {
		return err
	}

	todo, err := t.service.GetTodo(id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, todo)
}

func (t *todoHandler) Post(ctx echo.Context) error {
	var todo Todo

	if err := ctx.Bind(&todo); err != nil {
		return err
	}

	if err := t.service.InsertTodo(&todo); err != nil {
		return err
	}

	return ctx.JSON(200, todo)
}
