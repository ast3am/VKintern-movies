package service

import (
	"context"
	"errors"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/google/uuid"
)

func (s *Service) CreateActor(ctx context.Context, actor models.Actor) error {
	if actor.Name == "" {
		return errors.New("empty name")
	}
	id, _ := uuid.NewUUID()
	err := s.db.CreateActor(ctx, id, actor)
	return err
}

func (s *Service) DeleteActor(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = s.db.DeleteActor(ctx, uid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateActor(ctx context.Context, id string, actor models.Actor) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	nowActor, err := s.db.GetActorByUUID(ctx, uid)
	if actor.Name == "" {
		actor.Name = nowActor.Name
	}
	if actor.Gender == "" {
		actor.Gender = nowActor.Gender
	}
	if actor.BirthDate.IsZero() {
		actor.BirthDate = nowActor.BirthDate
	}
	err = s.db.UpdateActor(ctx, uid, actor)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetActorList(ctx context.Context) (map[string][]string, error) {
	result, err := s.db.GetActorList(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}
