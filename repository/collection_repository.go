package repository

import (
	"context"
	"fmt"

	"github.com/anilsenay/go-htmx-example/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CollectionRepository struct {
	database *pgxpool.Pool
}

func NewCollectionRepository(db *pgxpool.Pool) *CollectionRepository {
	return &CollectionRepository{
		database: db,
	}
}

func (r *CollectionRepository) GetCollection(id int) (model.Collection, error) {
	row := r.database.QueryRow(context.Background(), "SELECT * FROM collection WHERE id = $1", id)
	collection := model.Collection{}
	err := row.Scan(
		&collection.Id,
		&collection.Name,
		&collection.Color,
	)
	if err != nil {
		return model.Collection{}, err
	}
	return collection, nil
}

func (r *CollectionRepository) GetCollections() ([]model.Collection, error) {
	rows, err := r.database.Query(context.Background(), "SELECT * FROM collection")
	if err != nil {
		return nil, err
	}
	collections, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Collection])
	if err != nil {
		return nil, err
	}
	return collections, nil
}

func (r *CollectionRepository) Insert(collection model.Collection) (model.Collection, error) {
	row := r.database.QueryRow(context.Background(), "INSERT INTO collection (name, color) VALUES ($1,$2) RETURNING id", collection.Name, collection.Color)

	id := 0
	err := row.Scan(&id)
	if err != nil {
		return model.Collection{}, err
	}

	collection.Id = id
	return collection, nil
}

func (r *CollectionRepository) Update(id int, updates model.Collection) (model.Collection, error) {
	row := r.database.QueryRow(context.Background(), "UPDATE collection SET name = $2, color = $3 WHERE id = $1 RETURNING id, name, color", id, updates.Name, updates.Color)

	collection := model.Collection{}
	err := row.Scan(&collection.Id, &collection.Name, &collection.Color)
	if err != nil {
		return model.Collection{}, err
	}

	return collection, nil
}

func (r *CollectionRepository) Delete(id int) error {
	res, err := r.database.Exec(context.Background(), "DELETE FROM collection WHERE id = $1", id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("collection with id %d could not be deleted", id)
	}
	return nil
}
