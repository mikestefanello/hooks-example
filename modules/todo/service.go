package todo

import (
	"errors"
	"sync"

	"github.com/samber/do"
)

func NewTodoService(i *do.Injector) (Service, error) {
	return &todoService{
		todos: make([]Todo, 0),
		mu:    sync.RWMutex{},
	}, nil
}

func (t *todoService) GetTodo(id int) (Todo, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if len(t.todos) < id+1 {
		return Todo{}, errors.New("not found")
	}
	return t.todos[id], nil
}

func (t *todoService) GetTodos() ([]Todo, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.todos, nil
}

func (t *todoService) InsertTodo(todo *Todo) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	todo.ID = len(t.todos)
	HookTodoPreInsert.Dispatch(todo)

	t.todos = append(t.todos, *todo)
	HookTodoInsert.Dispatch(*todo)

	return nil
}
