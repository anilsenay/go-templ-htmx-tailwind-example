package repository

import (
	"fmt"
	"slices"
	"sync"

	"github.com/anilsenay/go-htmx-example/model"
)

var initialData = []model.Todo{
	{Id: 1, Text: "Todo 1", Done: true},
	{Id: 2, Text: "Todo 2"},
}

type TodoRepository struct {
	todos  []model.Todo
	mutex  sync.RWMutex
	nextId int
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		todos:  initialData,
		nextId: len(initialData) + 1,
	}
}

func (r *TodoRepository) RetrieveAll() ([]model.Todo, error) {
	return r.todos, nil
}

func (r *TodoRepository) Insert(todo model.Todo) (model.Todo, error) {
	r.mutex.Lock()
	todo.Id = r.nextId
	r.nextId++
	r.todos = append(r.todos, todo)
	r.mutex.Unlock()
	return todo, nil
}

func (r *TodoRepository) ChangeDone(id int) (*model.Todo, error) {
	idx := slices.IndexFunc(r.todos, func(e model.Todo) bool {
		return e.Id == id
	})

	if idx == -1 {
		return nil, fmt.Errorf("todo with id:%d not found", id)
	}
	r.mutex.Lock()
	r.todos[idx].Done = !r.todos[idx].Done
	r.mutex.Unlock()
	return &r.todos[idx], nil
}

func (r *TodoRepository) Delete(id int) error {
	idx := slices.IndexFunc(r.todos, func(e model.Todo) bool {
		return e.Id == id
	})
	if idx == -1 {
		return fmt.Errorf("todo with id:%d not found", id)
	}

	r.mutex.Lock()
	r.todos = append(r.todos[:idx], r.todos[idx+1:]...)
	r.mutex.Unlock()
	return nil
}
