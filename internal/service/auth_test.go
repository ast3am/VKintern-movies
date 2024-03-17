package service

import (
	"context"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/ast3am/VKintern-movies/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Auth(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	ctx := context.TODO()
	s := NewService(mockDB, log)
	DBexpected := &models.User{
		Email:    "admin@vk.ru",
		Password: "adminPassword#1",
		Role:     "admin",
	}
	email := "admin@vk.ru"
	password := "adminPassword#1"
	mockDB.On("GetUserByEmail", ctx, email).Return(DBexpected, nil)
	//Данные верны
	token, err := s.Auth(ctx, email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	//Неверный пароль
	password = "WrongPassword"
	token, err = s.Auth(ctx, email, password)
	assert.EqualError(t, err, "wrong email or password")
	assert.Equal(t, "", token)

}

func TestService_CheckToken(t *testing.T) {
	mockDB := mocks.NewDb(t)
	log := mocks.NewLogger(t)
	s := NewService(mockDB, log)
	mockPermisionLevel := "admin"
	userTokenNovalid := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQG1haWwuY29tIiwiZXhwIjoxNzEwNjE1OTUyLCJyb2xlIjoidXNlciJ9.3yM4fsdPvQUY5JgAIHRYIRwlFYs6Q2bBk0EaN_E7RoY"
	adminTokenValid := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQHZrLnJ1IiwiZXhwIjoxNzE1ODQ5MDg4LCJyb2xlIjoiYWRtaW4ifQ.oHRKkhFxnRcxE6etJ8VdeL4x9M2tiAUpxrOY1Bmluuc"
	//CheckToken true
	err := s.CheckToken(adminTokenValid, mockPermisionLevel)
	assert.NoError(t, err)
	//CheckToken false
	err = s.CheckToken(userTokenNovalid, mockPermisionLevel)
	assert.NotNil(t, err)
}
