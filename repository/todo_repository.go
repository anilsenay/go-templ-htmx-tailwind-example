package repository

import (
	"context"
	"fmt"

	"github.com/anilsenay/go-htmx-example/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository struct {
	database *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{
		database: db,
	}
}

func (r *TodoRepository) GetTodoByCollection(collection model.Collection) (model.CollectionWithTodoList, error) {
	rows, err := r.database.Query(context.Background(), "SELECT id, text, done FROM todo WHERE collection_id = $1", collection.Id)
	if err != nil {
		return model.CollectionWithTodoList{}, err
	}

	todo, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Todo])
	if err != nil {
		return model.CollectionWithTodoList{}, err
	}

	return model.CollectionWithTodoList{Collection: collection, List: todo}, nil
}

func (r *TodoRepository) Insert(collectionId int, todo model.Todo) (model.Todo, error) {
	row := r.database.QueryRow(context.Background(), "INSERT INTO todo (collection_id, text) VALUES ($1,$2) RETURNING id", collectionId, todo.Text)

	id := 0
	err := row.Scan(&id)
	if err != nil {
		return model.Todo{}, err
	}

	todo.Id = id
	return todo, nil
}

func (r *TodoRepository) ChangeDone(collectionId int, id int) (*model.Todo, error) {
	row := r.database.QueryRow(context.Background(), "UPDATE todo SET done = NOT done WHERE id = $1 RETURNING id, text, done", id)

	todo := model.Todo{}
	err := row.Scan(&todo.Id, &todo.Text, &todo.Done)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepository) Delete(collectionId int, id int) error {
	res, err := r.database.Exec(context.Background(), "DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("todo with id %d could not be deleted", id)
	}
	return nil
}
