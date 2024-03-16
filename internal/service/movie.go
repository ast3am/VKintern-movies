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
