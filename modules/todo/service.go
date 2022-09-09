package todo

import (
	"errors"

	"github.com/mikestefanello/hooks-example/services/cache"
	"github.com/samber/do"
)

const cacheKey = "todos"

func NewTodoService(i *do.Injector) (Service, error) {
	return &todoService{
		cache: do.MustInvoke[cache.Cache](i),
	}, nil
}

func (t *todoService) GetTodo(id int) (Todo, error) {
	todos := t.load()

	if len(todos) < id+1 {
		return Todo{}, errors.New("not found")
	}

	return todos[id], nil
}

func (t *todoService) GetTodos() ([]Todo, error) {
	return t.load(), nil
}

func (t *todoService) InsertTodo(todo *Todo) error {
	todos := t.load()

	todo.ID = len(todos)
	HookTodoPreInsert.Dispatch(todo)

	todos = append(todos, *todo)
	if err := t.save(todos); err != nil {
		return err
	}

	HookTodoInsert.Dispatch(*todo)

	return nil
}

func (t *todoService) load() []Todo {
	data, err := t.cache.Get(cacheKey)
	if err != nil {
		return make([]Todo, 0)
	}

	return data.([]Todo)
}

func (t *todoService) save(todos []Todo) error {
	return t.cache.Set(cacheKey, todos)
}
