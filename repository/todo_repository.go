package repository

import (
	"fmt"
	"slices"
	"sync"

	"github.com/anilsenay/go-htmx-example/model"
)

var initialCollections = []model.Collection{
	{Id: 1, Name: "My Todo", HexColor: "#d489d9"},
	{Id: 2, Name: "My Todo 2", HexColor: "#86dbdb"},
}

var todoByCollection = map[int][]model.Todo{
	1: {{Id: 1, Text: "List 1 Todo 1", Done: true}, {Id: 2, Text: "List 1 Todo 2"}},
	2: {{Id: 3, Text: "List 2 Todo 1", Done: true}, {Id: 4, Text: "List 2 Todo 2"}},
}

type TodoRepository struct {
	collections      []model.Collection
	todoByCollection map[int][]model.Todo
	mutex            sync.RWMutex
	nextTodoId       int
	nextCollectionId int
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		collections:      initialCollections,
		todoByCollection: todoByCollection,
		nextTodoId:       5,
		nextCollectionId: 3,
	}
}

func (r *TodoRepository) findCollectionById(id int) *model.Collection {
	idx := slices.IndexFunc(r.collections, func(e model.Collection) bool {
		return e.Id == id
	})
	if idx == -1 {
		return nil
	}
	return &r.collections[idx]
}

func (r *TodoRepository) findTodoById(collectionId int, id int) *model.Todo {
	todo, found := r.todoByCollection[collectionId]
	if !found {
		return nil
	}
	idx := slices.IndexFunc(todo, func(e model.Todo) bool {
		return e.Id == id
	})
	if idx == -1 {
		return nil
	}
	return &todo[idx]
}

func (r *TodoRepository) GetCollections() ([]model.Collection, error) {
	return r.collections, nil
}

func (r *TodoRepository) GetTodoByCollection(id int) (model.CollectionWithTodoList, error) {
	collection := r.findCollectionById(id)
	if collection == nil {
		return model.CollectionWithTodoList{}, fmt.Errorf("list with id %d not found", id)
	}
	return model.CollectionWithTodoList{Collection: *collection, List: r.todoByCollection[id]}, nil
}

func (r *TodoRepository) Insert(collectionId int, todo model.Todo) (model.Todo, error) {
	_, found := r.todoByCollection[collectionId]
	if !found {
		return todo, fmt.Errorf("list with id %d not found", collectionId)
	}

	r.mutex.Lock()
	todo.Id = r.nextTodoId
	r.nextTodoId++
	r.todoByCollection[collectionId] = append(r.todoByCollection[collectionId], todo)
	r.mutex.Unlock()
	return todo, nil
}

func (r *TodoRepository) ChangeDone(collectionId int, id int) (*model.Todo, error) {
	todo := r.findTodoById(collectionId, id)
	if todo == nil {
		return nil, fmt.Errorf("todo with id:%d not found", id)
	}

	r.mutex.Lock()
	todo.Done = !todo.Done
	r.mutex.Unlock()
	return todo, nil
}

func (r *TodoRepository) Delete(collectionId int, id int) error {
	todo, found := r.todoByCollection[collectionId]
	if !found {
		return fmt.Errorf("collection with id:%d not found", collectionId)
	}

	idx := slices.IndexFunc(todo, func(e model.Todo) bool {
		return e.Id == id
	})
	if idx == -1 {
		return fmt.Errorf("todo with id:%d not found", id)
	}

	r.mutex.Lock()
	r.todoByCollection[collectionId] = append(r.todoByCollection[collectionId][:idx], r.todoByCollection[collectionId][idx+1:]...)
	r.mutex.Unlock()
	return nil
}

func (r *TodoRepository) InsertCollection(collection model.Collection) (model.Collection, error) {
	r.mutex.Lock()
	collection.Id = r.nextCollectionId
	r.nextCollectionId++
	r.collections = append(r.collections, collection)
	r.todoByCollection[collection.Id] = []model.Todo{}
	r.mutex.Unlock()
	return collection, nil
}
