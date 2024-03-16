package db

import (
	"context"
	"errors"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (db *DB) CreateMovie(ctx context.Context, id uuid.UUID, movie models.Movie) error {
	tx, err := db.dbConnect.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	createOrder := `
	INSERT INTO movies values ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, createOrder, id, movie.Name, movie.Description, movie.ReleaseDate, movie.Rating)
	if err != nil {
		return err
	}
	for _, val := range movie.ActorList {
		updateFilmOrder := `
		INSERT INTO movie_actors (movie_uuid, actor_name)
		VALUES ($1, $2)`
		_, err = tx.Exec(ctx, updateFilmOrder, id, val)
		if err != nil {
			return err
		}
		var actorID string
		selectActorOrder := `
		SELECT uuid FROM actors WHERE actors.name = $1`
		err = tx.QueryRow(ctx, selectActorOrder, val).Scan(&actorID)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		} else if err != nil {
			return err
		}
		actorUID, _ := uuid.Parse(actorID)
		updateOrder := `
		UPDATE movie_actors
		SET actor_uuid = $1
		WHERE movie_uuid = $2 AND actor_name = $3`
		_, err = tx.Exec(ctx, updateOrder, actorUID, id, val)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
