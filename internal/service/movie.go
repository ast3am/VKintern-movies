package service

import (
	"context"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/google/uuid"
)

func (s *Service) CreateMovie(ctx context.Context, movie models.Movie) error {
	err := movie.Validate()
	if err != nil {
		return err
	}
	id, _ := uuid.NewUUID()
	err = s.db.CreateMovie(ctx, id, movie)
	return err
}

func (s *Service) UpdateMovie(ctx context.Context, id string, movie models.Movie) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	nowMovie, err := s.db.GetMovieByUUID(ctx, uid)
	if err != nil {
		return err
	}
	if movie.Name == "" {
		movie.Name = nowMovie.Name
	}
	if movie.Description == "" {
		movie.Description = nowMovie.Description
	}
	if movie.ReleaseDate.IsZero() {
		movie.ReleaseDate = nowMovie.ReleaseDate
	}
	if movie.Rating == 0.0 {
		movie.Rating = nowMovie.Rating
	}
	if len(movie.ActorList) == 0 {
		movie.ActorList = nowMovie.ActorList
	}
	err = movie.Validate()
	if err != nil {
		return err
	}
	err = s.db.UpdateMovie(ctx, uid, movie)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteMovie(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = s.db.DeleteMovie(ctx, uid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetMovieList(ctx context.Context, sortby, line string) ([]*models.Movie, error) {
	if sortby == "" {
		sortby = "rating"
	}
	if line == "" {
		line = "desc"
	}
	result, err := s.db.GetMovieList(ctx, sortby, line)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) GetMovie(ctx context.Context, actor, movie string) ([]*models.Movie, error) {
	actor = "%" + actor + "%"
	movie = "%" + movie + "%"
	result, err := s.db.GetMovie(ctx, actor, movie)
	if err != nil {
		return nil, err
	}
	return result, nil
}
