package handlers

import (
	"bytes"
	"errors"
	"github.com/ast3am/VKintern-movies/api/handlers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Auth(t *testing.T) {
	mux := http.NewServeMux()
	serv := mocks.NewService(t)
	log := mocks.NewLogger(t)
	h := NewHandler(serv, log)
	h.RegisterHandlers(mux)

	testTable := []struct {
		name                 string
		requestBody          []byte
		httpMethod           string
		expectedStatusCode   int
		expectedResponseBody []byte
	}{
		{
			"positive",
			[]byte(`{"email": "test@example.com", "password": "password123"}`),
			http.MethodPost,
			http.StatusOK,
			[]byte(`authorization success{"token":"some token"}` + "\n"),
		}, {
			"wrong method",
			[]byte(`{"email": "test@example.com", "password": "password123"}`),
			http.MethodGet,
			http.StatusMethodNotAllowed,
			[]byte(`invalid request method` + "\n"),
		}, {
			"wrong json",
			[]byte(`{wrong json}`),
			http.MethodPost,
			http.StatusUnprocessableEntity,
			[]byte(`JSON parsing error` + "\n"),
		}, {
			"wrong user",
			[]byte(`{"email": "test@example.com", "password": "password123"}`),
			http.MethodPost,
			http.StatusBadRequest,
			[]byte(`wrong email or password` + "\n"),
		},
	}
	for _, test := range testTable {
		if test.name == "positive" {
			serv.On("Auth", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("some token", nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, "authorization").Return(0)
		} else if test.name == "wrong method" {
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(MethodNotAllowed)).Return(0)
		} else if test.name == "wrong json" {
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(ParsingJSONError)).Return(0)
		} else if test.name == "wrong user" {
			serv.On("Auth", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("", errors.New("wrong email or password")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("wrong email or password")).Return(0)
		}
		req, err := http.NewRequest(test.httpMethod, "/auth", bytes.NewBuffer(test.requestBody))
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, req)
		responseBody := r.Body.Bytes()

		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, r.Code)
		assert.Equal(t, test.expectedResponseBody, responseBody)
	}
}
