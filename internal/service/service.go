package service

import (
	"context"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/google/uuid"
)

type db interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetActorByUUID(ctx context.Context, id uuid.UUID) (*models.Actor, error)
	CreateActor(ctx context.Context, id uuid.UUID, actor models.Actor) error
	UpdateActor(ctx context.Context, id uuid.UUID, actor models.Actor) error
	GetActorList(ctx context.Context) (map[string][]string, error)
	DeleteActor(ctx context.Context, uid uuid.UUID) error
	CreateMovie(ctx context.Context, id uuid.UUID, movie models.Movie) error
}

type logger interface {
	DebugMsg(msg string)
	ErrorMsg(msg string, err error)
}

type Service struct {
	db  db
	log logger
}

func NewService(db db, log logger) *Service {
	return &Service{
		db,
		log,
	}
}
