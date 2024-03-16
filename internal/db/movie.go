package db

import (
	"context"
	"database/sql"
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

func (db *DB) GetMovieByUUID(ctx context.Context, id uuid.UUID) (*models.Movie, error) {
	result := models.Movie{}
	queryOrder := `
	SELECT name, description, release_date, rating FROM movies WHERE uuid = $1
	`
	err := db.dbConnect.QueryRow(ctx, queryOrder, id).Scan(&result.Name, &result.Description, &result.ReleaseDate, &result.Rating)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (db *DB) UpdateMovie(ctx context.Context, id uuid.UUID, movie models.Movie) error {
	tx, err := db.dbConnect.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updateMoviesOrder := `
	UPDATE movies
	SET name = $2, description = $3, release_date = $4, rating = $5
	WHERE uuid = $1
	`
	_, err = tx.Exec(ctx, updateMoviesOrder, id, movie.Name, movie.Description, movie.ReleaseDate, movie.Rating)
	if err != nil {
		return err
	}

	if len(movie.ActorList) != 0 {
		deleteFilmOrder := `
		DELETE FROM movie_actors WHERE movie_uuid = $1`
		_, err = tx.Exec(ctx, deleteFilmOrder, id)
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
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteMovie(ctx context.Context, uid uuid.UUID) error {
	tx, err := db.dbConnect.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	deleteOrder := `
	DELETE FROM movies WHERE uuid = $1
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

func (db *DB) GetMovieList(ctx context.Context, sortby, list string) ([]*models.Movie, error) {
	result := make([]*models.Movie, 0)
	getMovieOrder := `SELECT name, description, release_date, rating, actor_name FROM movies
	JOIN movie_actors ma ON ma.movie_uuid = movies.uuid`
	getMovieOrder += " order by " + sortby + " " + list
	rows, err := db.dbConnect.Query(ctx, getMovieOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var timestamp sql.NullTime
	movieName := ""
	i := 0
	for rows.Next() {
		actorName := ""
		movieStruct := models.Movie{}
		err = rows.Scan(&movieStruct.Name, &movieStruct.Description, &timestamp, &movieStruct.Rating, &actorName)
		if err != nil {
			return nil, err
		}
		if timestamp.Valid {
			movieStruct.ReleaseDate = timestamp.Time
		}
		if movieStruct.Name != movieName {
			movieName = movieStruct.Name
			movieStruct.ActorList = append(movieStruct.ActorList, actorName)
			result = append(result, &movieStruct)
			i++
			continue
		} else {
			result[i-1].ActorList = append(result[i-1].ActorList, actorName)
		}
	}
	return result, nil
}

func (db *DB) GetMovie(ctx context.Context, actor, movie string) ([]*models.Movie, error) {
	result := make([]*models.Movie, 0)
	getMovieOrder := `
	SELECT name, description, release_date, rating, actor_name FROM movies
	JOIN movie_actors ma ON ma.movie_uuid = movies.uuid
	WHERE movies.name LIKE $1 AND ma.actor_name LIKE $2
	order by name`
	rows, err := db.dbConnect.Query(ctx, getMovieOrder, movie, actor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var timestamp sql.NullTime
	movieName := ""
	i := 0
	for rows.Next() {
		actorName := ""
		movieStruct := models.Movie{}
		err = rows.Scan(&movieStruct.Name, &movieStruct.Description, &timestamp, &movieStruct.Rating, &actorName)
		if err != nil {
			return nil, err
		}
		if timestamp.Valid {
			movieStruct.ReleaseDate = timestamp.Time
		}
		if movieStruct.Name != movieName {
			movieName = movieStruct.Name
			movieStruct.ActorList = append(movieStruct.ActorList, actorName)
			result = append(result, &movieStruct)
			i++
			continue
		} else {
			result[i-1].ActorList = append(result[i-1].ActorList, actorName)
		}
	}
	return result, nil
}
