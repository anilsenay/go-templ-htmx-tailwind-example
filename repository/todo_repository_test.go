package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/anilsenay/go-htmx-example/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestTodoRepository_Insert(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "123456", "go-htmx-example")

	cfg, _ := pgxpool.ParseConfig(psqlInfo)
	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	todoRepository := NewTodoRepository(db)
	collectionRepository := NewCollectionRepository(db)

	var collectionId int
	t.Run("Insert Collection", func(t *testing.T) {
		c, err := collectionRepository.Insert(model.Collection{Name: "Collection1", Color: "#aaaaaa"})
		assert.NoError(t, err)
		assert.NotEmpty(t, c.Id)
		collectionId = c.Id
	})

	t.Run("Get Collection", func(t *testing.T) {
		collection, err := collectionRepository.GetCollection(collectionId)
		assert.NoError(t, err)
		assert.NotEmpty(t, collection.Id)
	})

	t.Run("Get Collections", func(t *testing.T) {
		collections, err := collectionRepository.GetCollections()
		assert.NoError(t, err)
		assert.NotEmpty(t, collections)
	})

	t.Run("Update Collection", func(t *testing.T) {
		collection, err := collectionRepository.Update(collectionId, model.Collection{Name: "Collection", Color: "#111111"})
		assert.NoError(t, err)
		assert.NotEmpty(t, collection.Id)
		assert.Equal(t, "Collection", collection.Name)
		assert.Equal(t, "#111111", collection.Color)
	})

	var todoId int
	t.Run("Insert Todo", func(t *testing.T) {
		todo, err := todoRepository.Insert(collectionId, model.Todo{Text: "test"})
		assert.NoError(t, err)
		assert.NotEmpty(t, todo.Id)
		todoId = todo.Id
	})

	t.Run("Get Todo By Collection", func(t *testing.T) {
		todo, err := todoRepository.GetTodoByCollection(model.Collection{Id: collectionId})
		assert.NoError(t, err)
		assert.NotEmpty(t, todo.Id)
		assert.Equal(t, todoId, todo.Id)
		assert.Equal(t, collectionId, todo.Collection.Id)
	})

	t.Run("Change todo", func(t *testing.T) {
		todo, err := todoRepository.ChangeDone(0, todoId)
		assert.NoError(t, err)
		assert.NotEmpty(t, todo)
		assert.Equal(t, true, todo.Done)

		// set false back
		todo, err = todoRepository.ChangeDone(0, todoId)
		assert.NoError(t, err)
		assert.NotEmpty(t, todo)
		assert.Equal(t, false, todo.Done)
	})

	t.Run("Delete Todo", func(t *testing.T) {
		err := todoRepository.Delete(0, todoId)
		assert.NoError(t, err)
	})

	t.Run("Delete Collection", func(t *testing.T) {
		err := collectionRepository.Delete(collectionId)
		assert.NoError(t, err)
	})
}
