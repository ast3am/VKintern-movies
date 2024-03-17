package db

import (
	"context"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/google/uuid"
)

func (db *DB) CreateActor(ctx context.Context, id uuid.UUID, actor models.Actor) error {
	tx, err := db.dbConnect.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	createOrder := `
	INSERT INTO actors values ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, createOrder, id, actor.Name, actor.Gender, actor.BirthDate)
	if err != nil {
		return err
	}
	updateFilmOrder := `
	UPDATE movie_actors
	SET actor_uuid = $1
	where actor_name = $2;
`
	_, err = tx.Exec(ctx, updateFilmOrder, id, actor.Name)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteActor(ctx context.Context, uid uuid.UUID) error {
	tx, err := db.dbConnect.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	deleteOrder := `
	DELETE FROM actors WHERE uuid = $1
	`
	_, err = tx.Exec(ctx, deleteOrder, uid)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetActorByUUID(ctx context.Context, id uuid.UUID) (*models.Actor, error) {
	result := models.Actor{}
	queryOrder := `
	SELECT name, gender, birth_date FROM actors WHERE uuid = $1
	`
	err := db.dbConnect.QueryRow(ctx, queryOrder, id).Scan(&result.Name, &result.Gender, &result.BirthDate)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (db *DB) UpdateActor(ctx context.Context, id uuid.UUID, actor models.Actor) error {
	tx, err := db.dbConnect.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updateOrder := `
	UPDATE actors
	SET name = $2, gender = $3, birth_date = $4
	WHERE uuid = $1
	`
	_, err = tx.Exec(ctx, updateOrder, id, actor.Name, actor.Gender, actor.BirthDate)
	if err != nil {
		return err
	}
	updateFilmOrder := `
	UPDATE movie_actors
	SET actor_name = $2
	where actor_uuid = $1;
`
	_, err = tx.Exec(ctx, updateFilmOrder, id, actor.Name)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetActorList(ctx context.Context) (map[string][]string, error) {
	result := make(map[string][]string)
	getActorsOrder := `
	SELECT actors.name AS actor_name, movies.name AS movie_name
	FROM movies
	JOIN movie_actors ON movies.uuid = movie_actors.movie_uuid
	JOIN actors ON movie_actors.actor_uuid = actors.uuid
	ORDER BY actor_name, movie_name`
	rows, err := db.dbConnect.Query(ctx, getActorsOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var actorName, film string
		err = rows.Scan(&actorName, &film)
		if err != nil {
			return nil, err
		}
		result[actorName] = append(result[actorName], film)
	}
	return result, nil
}
