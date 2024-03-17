package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	noValidToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQG1haWwuY29tIiwiZXhwIjoxNzEwNjE1OTUyLCJyb2xlIjoidXNlciJ9.3yM4fsdPvQUY5JgAIHRYIRwlFYs6Q2bBk0EaN_E7RoY"
	validAdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQHZrLnJ1IiwiZXhwIjoxNzE1ODQ5MDg4LCJyb2xlIjoiYWRtaW4ifQ.oHRKkhFxnRcxE6etJ8VdeL4x9M2tiAUpxrOY1Bmluuc"
	validUserToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyQG1haWwuY29tIiwiZXhwIjoxNzE1ODQ5MTI1LCJyb2xlIjoidXNlciJ9.6iQm2423g4AurdEMpl5tMZ3a-75OucUO17LqWkgPCak"
)

func TestGetToken(t *testing.T) {
	testMail := "test@mail.ru"
	testRole := "user"
	token, err := GetToken(testMail, testRole)
	assert.NotNil(t, token)
	assert.Nil(t, err)
}

func TestCheckPermissionByToken(t *testing.T) {
	err := CheckPermissionByToken(validAdminToken, "admin")
	assert.Nil(t, err)
	err = CheckPermissionByToken(noValidToken, "User")
	assert.EqualError(t, err, "not a valid token")
	err = CheckPermissionByToken(validUserToken, "admin")
	assert.EqualError(t, err, "permission denied")
}
