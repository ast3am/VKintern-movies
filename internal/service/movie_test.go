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
)

func TestService_CreateMovie(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	movieTrue := models.Movie{
		Name:        "Forrest Gump",
		Description: "Description 1",
		ReleaseDate: GetDate("1994-07-06"),
		Rating:      8.8,
		ActorList: []string{
			"Tom Hanks",
			"Brad Pitt",
		},
	}
	movieFalse := models.Movie{
		Name:        "",
		Description: "Description 1",
		ReleaseDate: GetDate("1994-07-06"),
		Rating:      8.8,
		ActorList: []string{
			"Tom Hanks",
			"Brad Pitt",
		},
	}

	mockDB.On("CreateMovie", ctx, mock.AnythingOfType("uuid.UUID"), movieTrue).Return(nil)
	//true
	err := s.CreateMovie(ctx, movieTrue)
	assert.Nil(t, err)
	//false
	err = s.CreateMovie(ctx, movieFalse)
	assert.NotNil(t, err)
}

func TestService_UpdateMovie(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	validID := "b0482c7a-1a4c-4a3c-9463-35f0036a0d60"
	validIDnoMovie := "7cfd7e7f-3a27-46c5-9c3e-3b3d3b1c1416"
	noValidID := "b0482c7a-1a4c-4a3"
	movie := &models.Movie{
		Name:        "Forrest Gump",
		Description: "Description 1",
		ReleaseDate: GetDate("1994-07-06"),
		Rating:      8.8,
		ActorList: []string{
			"Tom Hanks",
			"Brad Pitt",
		},
	}
	mockDB.On("GetMovieByUUID", ctx, mock.AnythingOfType("uuid.UUID")).Return(movie, func(ctx context.Context, id uuid.UUID) error {
		if id.String() == validID {
			return nil
		}
		return errors.New("some error")
	})
	mockDB.On("UpdateMovie", ctx, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("models.Movie")).Return(nil)
	testData := []models.Movie{
		{
			Name:        "Training Day",
			Description: "Description 2",
			ReleaseDate: GetDate("1999-07-06"),
			Rating:      9.8,
			ActorList: []string{
				"Tom Hanks",
				"Brad Pitt",
			},
		}, {
			Name:        "",
			Description: "",
			ReleaseDate: GetDate("1999-07-06"),
			ActorList: []string{
				"Tom Hanks",
				"Brad Pitt",
			},
		}, {
			Name:        "Training Day",
			Description: "Description 2",
			Rating:      9.8,
		},
	}
	for _, val := range testData {
		err := s.UpdateMovie(ctx, validID, val)
		assert.Nil(t, err)
	}
	//novalid UUID
	err := s.UpdateMovie(ctx, noValidID, testData[0])
	assert.NotNil(t, err)
	//noUser
	err = s.UpdateMovie(ctx, validIDnoMovie, testData[0])
	assert.NotNil(t, err)

	// no valid movie
	noValidMovie := models.Movie{
		Name:        "Forrest Gump",
		Description: "Description 1",
		ReleaseDate: GetDate("1994-07-06"),
		Rating:      12,
		ActorList: []string{
			"Tom Hanks",
			"Brad Pitt",
		},
	}
	err = s.UpdateMovie(ctx, validID, noValidMovie)
	assert.NotNil(t, err)
}

func TestService_DeleteMovie(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	validID := "b0482c7a-1a4c-4a3c-9463-35f0036a0d60"
	validIDnoUser := "7cfd7e7f-3a27-46c5-9c3e-3b3d3b1c1416"
	noValidID := "b0482c7a-1a4c-4a3"
	mockDB.On("DeleteMovie", ctx, mock.AnythingOfType("uuid.UUID")).Return(func(ctx context.Context, id uuid.UUID) error {
		if id.String() == validID {
			return nil
		}
		return errors.New("some error")
	})
	//true
	err := s.DeleteMovie(ctx, validID)
	assert.Nil(t, err)
	//novalid UUID
	err = s.DeleteMovie(ctx, noValidID)
	assert.NotNil(t, err)
	//noUser
	err = s.DeleteMovie(ctx, validIDnoUser)
	assert.NotNil(t, err)
}

func TestService_GetMovieList(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	result := []*models.Movie{
		{
			Name:        "Forrest Gump",
			Description: "Description 1",
			ReleaseDate: GetDate("1994-07-06"),
			Rating:      8.5,
			ActorList: []string{
				"Tom Hanks",
				"Brad Pitt",
			},
		}, {
			Name:        "Training Day",
			Description: "Description 2",
			ReleaseDate: GetDate("2001-10-05"),
			Rating:      7.7,
			ActorList: []string{
				"Brad Pitt",
				"Cate Blanchett",
			},
		},
	}
	mockDB.On("GetMovieList", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(result, nil).Times(2)
	result, err := s.GetMovieList(ctx, "", "")
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = s.GetMovieList(ctx, "rating", "desc")
	assert.NotNil(t, result)
	assert.Nil(t, err)
	mockDB.On("GetMovieList", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("some error")).Once()
	result, err = s.GetMovieList(ctx, "rating", "desc")
	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestService_GetMovie(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	result := []*models.Movie{
		{
			Name:        "Forrest Gump",
			Description: "Description 1",
			ReleaseDate: GetDate("1994-07-06"),
			Rating:      8.5,
			ActorList: []string{
				"Tom Hanks",
				"Brad Pitt",
			},
		}, {
			Name:        "Training Day",
			Description: "Description 2",
			ReleaseDate: GetDate("2001-10-05"),
			Rating:      7.7,
			ActorList: []string{
				"Brad Pitt",
				"Cate Blanchett",
			},
		},
	}
	mockDB.On("GetMovie", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(result, nil).Once()
	result, err := s.GetMovie(ctx, "for", "bra")
	assert.NotNil(t, result)
	assert.Nil(t, err)
	mockDB.On("GetMovie", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("some error")).Once()
	result, err = s.GetMovie(ctx, "for", "bra")
	assert.Nil(t, result)
	assert.NotNil(t, err)
}
