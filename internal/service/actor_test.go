package service

import (
	"context"
	"errors"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/ast3am/VKintern-movies/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func GetDate(date string) time.Time {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, date)
	return t
}

func TestService_CreateActor(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	actorTrue := models.Actor{
		Name:      "Tom Hanks",
		Gender:    "Male",
		BirthDate: GetDate("1956-07-09"),
	}
	actorFalse := models.Actor{
		Name:      "",
		Gender:    "Male",
		BirthDate: GetDate("1956-07-09"),
	}

	mockDB.On("CreateActor", ctx, mock.AnythingOfType("uuid.UUID"), actorTrue).Return(nil)
	//true
	err := s.CreateActor(ctx, actorTrue)
	assert.Nil(t, err)
	//false
	err = s.CreateActor(ctx, actorFalse)
	assert.EqualError(t, err, "empty name")
}

func TestService_DeleteActor(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	validID := "b0482c7a-1a4c-4a3c-9463-35f0036a0d60"
	validIDnoUser := "7cfd7e7f-3a27-46c5-9c3e-3b3d3b1c1416"
	noValidID := "b0482c7a-1a4c-4a3"
	mockDB.On("DeleteActor", ctx, mock.AnythingOfType("uuid.UUID")).Return(func(ctx context.Context, id uuid.UUID) error {
		if id.String() == validID {
			return nil
		}
		return errors.New("some error")
	})
	//true
	err := s.DeleteActor(ctx, validID)
	assert.Nil(t, err)
	//novalid UUID
	err = s.DeleteActor(ctx, noValidID)
	assert.NotNil(t, err)
	//noUser
	err = s.DeleteActor(ctx, validIDnoUser)
	assert.NotNil(t, err)
}

func TestService_UpdateActor(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	validID := "b0482c7a-1a4c-4a3c-9463-35f0036a0d60"
	validIDnoUser := "7cfd7e7f-3a27-46c5-9c3e-3b3d3b1c1416"
	noValidID := "b0482c7a-1a4c-4a3"
	actor := &models.Actor{
		Name:      "Tom Hanks",
		Gender:    "Male",
		BirthDate: GetDate("1956-07-09"),
	}
	mockDB.On("GetActorByUUID", ctx, mock.AnythingOfType("uuid.UUID")).Return(actor, func(ctx context.Context, id uuid.UUID) error {
		if id.String() == validID {
			return nil
		}
		return errors.New("some error")
	})
	mockDB.On("UpdateActor", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("models.Actor")).Return(nil)
	testData := []models.Actor{
		{
			Name:   "Tom Hanks",
			Gender: "Female",
		},
		{
			Name:      "",
			Gender:    "",
			BirthDate: GetDate("1956-07-09"),
		},
	}
	for _, val := range testData {
		err := s.UpdateActor(ctx, validID, val)
		assert.Nil(t, err)
	}
	//novalid UUID
	err := s.UpdateActor(ctx, noValidID, testData[0])
	assert.NotNil(t, err)
	//noUser
	err = s.UpdateActor(ctx, validIDnoUser, testData[0])
	assert.NotNil(t, err)
}

func TestService_GetActorList(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	mockDB.On("GetActorList", ctx).Return(map[string][]string{
		"Tom Hanks": []string{
			"Forrest Gump",
			"Training Day",
		},
	}, nil).Once()
	result, err := s.GetActorList(ctx)
	assert.NotNil(t, result)
	assert.Nil(t, err)

	mockDB.On("GetActorList", ctx).Return(nil, errors.New("some error")).Once()
	result, err = s.GetActorList(ctx)
	assert.Nil(t, result)
	assert.NotNil(t, err)
}
