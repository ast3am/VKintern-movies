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

func TestHandler_CreateActor(t *testing.T) {
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
			"wrong method",
			[]byte(`{"name": "Tom Hanks", "gender": "male", "birth_date": "2000-01-01"}`),
			http.MethodGet,
			http.StatusMethodNotAllowed,
			[]byte(`invalid request method` + "\n"),
		}, {
			"wrong token",
			[]byte(`{"name": "Tom Hanks", "gender": "male", "birth_date": "2000-01-01"}`),
			http.MethodPost,
			http.StatusUnauthorized,
			[]byte(`invalid token` + "\n"),
		}, {
			"wrong json",
			[]byte(`"some data"`),
			http.MethodPost,
			http.StatusUnprocessableEntity,
			[]byte(`JSON parsing error` + "\n"),
		}, {
			"no valid data",
			[]byte(`{"name": "", "gender": "male", "birth_date": "2000-01-01"}`),
			http.MethodPost,
			http.StatusUnprocessableEntity,
			[]byte(`no valid data` + "\n"),
		}, {
			"another error service",
			[]byte(`{"name": "Tom Hanks", "gender": "male", "birth_date": "2000-01-01"}`),
			http.MethodPost,
			http.StatusBadRequest,
			[]byte(`service error` + "\n"),
		}, {
			"positive",
			[]byte(`{"name": "Tom Hanks", "gender": "male", "birth_date": "2000-01-01"}`),
			http.MethodPost,
			http.StatusOK,
			[]byte(`Actor created`),
		},
	}
	for _, test := range testTable {
		if test.name == "wrong method" {
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(MethodNotAllowed)).Return(0)
		} else if test.name == "wrong token" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(InvalidToken)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(InvalidToken)).Return(0)
		} else if test.name == "wrong json" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), mock.AnythingOfType("*json.UnmarshalTypeError")).Return(0)
		} else if test.name == "no valid data" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("CreateActor", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("models.Actor")).Return(errors.New(UnprocessableEntity)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(UnprocessableEntity)).Return(0)
		} else if test.name == "another error service" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("CreateActor", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("models.Actor")).Return(errors.New("service error")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("service error")).Return(0)
		} else if test.name == "positive" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("CreateActor", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("models.Actor")).Return(nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string")).Return(0)
		}
		req, err := http.NewRequest(test.httpMethod, "/actor/create", bytes.NewBuffer(test.requestBody))
		req.Header.Set("Authorization", "Test-token")
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, req)
		responseBody := r.Body.Bytes()

		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, r.Code)
		assert.Equal(t, test.expectedResponseBody, responseBody)
	}
}

func TestHandler_GetActorsList(t *testing.T) {
	mux := http.NewServeMux()
	serv := mocks.NewService(t)
	log := mocks.NewLogger(t)
	h := NewHandler(serv, log)
	h.RegisterHandlers(mux)

	testTable := []struct {
		name                 string
		httpMethod           string
		expectedStatusCode   int
		expectedResponseBody []byte
	}{
		{
			"wrong method",
			http.MethodPost,
			http.StatusMethodNotAllowed,
			[]byte(`invalid request method` + "\n"),
		}, {
			"wrong token",
			http.MethodGet,
			http.StatusUnauthorized,
			[]byte(`invalid token` + "\n"),
		}, {
			"service problem",
			http.MethodGet,
			http.StatusBadRequest,
			[]byte(`service problem` + "\n"),
		}, {
			"positive",
			http.MethodGet,
			http.StatusOK,
			[]byte(`{"boris":["snatch"],"jhon":["badBoys2","badBoys"],"suize":["badBoys","hatiko"]}`),
		},
	}
	for _, test := range testTable {
		if test.name == "wrong method" {
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(MethodNotAllowed)).Return(0).Once()
		} else if test.name == "wrong token" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(InvalidToken)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(InvalidToken)).Return(0).Once()
		} else if test.name == "service problem" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
			serv.On("GetActorList", mock.AnythingOfType("context.backgroundCtx")).Return(nil, errors.New("service problem")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("service problem")).Return(0).Once()
		} else if test.name == "positive" {
			result := map[string][]string{"boris": {"snatch"}, "jhon": {"badBoys2", "badBoys"}, "suize": {"badBoys", "hatiko"}}
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
			serv.On("GetActorList", mock.AnythingOfType("context.backgroundCtx")).Return(result, nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string")).Return(0).Once()
		}
		req, err := http.NewRequest(test.httpMethod, "/actor/get-list", nil)
		req.Header.Set("Authorization", "Test-token")
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, req)
		responseBody := r.Body.Bytes()
		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, r.Code)
		assert.Equal(t, test.expectedResponseBody, responseBody)
	}
}

func TestHandler_DeleteActor(t *testing.T) {
	mux := http.NewServeMux()
	serv := mocks.NewService(t)
	log := mocks.NewLogger(t)
	h := NewHandler(serv, log)
	h.RegisterHandlers(mux)
	testTable := []struct {
		name                 string
		httpMethod           string
		id                   string
		expectedStatusCode   int
		expectedResponseBody []byte
	}{
		{
			"wrong method",
			http.MethodPost,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusMethodNotAllowed,
			[]byte(`invalid request method` + "\n"),
		}, {
			"wrong token",
			http.MethodDelete,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusUnauthorized,
			[]byte(`invalid token` + "\n"),
		}, {
			"service problem",
			http.MethodDelete,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusBadRequest,
			[]byte(`service problem` + "\n"),
		}, {
			"positive",
			http.MethodDelete,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusOK,
			[]byte(`Actor deleted`),
		},
	}
	for _, test := range testTable {
		if test.name == "wrong method" {
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(MethodNotAllowed)).Return(0).Once()
		} else if test.name == "wrong token" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(InvalidToken)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(InvalidToken)).Return(0).Once()
		} else if test.name == "service problem" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
			serv.On("DeleteActor", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).Return(errors.New("service problem")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("service problem")).Return(0).Once()
		} else if test.name == "positive" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
			serv.On("DeleteActor", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).Return(nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string")).Return(0).Once()
		}
		req, err := http.NewRequest(test.httpMethod, "/actor/delete/"+test.id, nil)
		req.Header.Set("Authorization", "Test-token")
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, req)
		responseBody := r.Body.Bytes()
		t.Logf("Extra Variable: %s", test.name)
		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, r.Code)
		assert.Equal(t, test.expectedResponseBody, responseBody)
	}
}

func TestHandler_UpdateActor(t *testing.T) {
	mux := http.NewServeMux()
	serv := mocks.NewService(t)
	log := mocks.NewLogger(t)
	h := NewHandler(serv, log)
	h.RegisterHandlers(mux)
	testTable := []struct {
		name                 string
		requestBody          []byte
		httpMethod           string
		id                   string
		expectedStatusCode   int
		expectedResponseBody []byte
	}{
		{
			"wrong method",
			[]byte(`{"name": "Tom Banks", "gender": "female", "birth_date": "2000-01-01"}`),
			http.MethodGet,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusMethodNotAllowed,
			[]byte(`invalid request method` + "\n"),
		}, {
			"wrong token",
			[]byte(`{"name": "Tom Banks", "gender": "female", "birth_date": "2000-01-01"}`),
			http.MethodPatch,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusUnauthorized,
			[]byte(`invalid token` + "\n"),
		}, {
			"wrong json",
			[]byte(`"some data"`),
			http.MethodPatch,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusUnprocessableEntity,
			[]byte(`JSON parsing error` + "\n"),
		}, {
			"service error",
			[]byte(`{"name": "Tom Banks", "gender": "female", "birth_date": "2000-01-01"}`),
			http.MethodPatch,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusBadRequest,
			[]byte(`service error` + "\n"),
		}, {
			"positive",
			[]byte(`{"name": "Tom Banks", "gender": "female", "birth_date": "2000-01-01"}`),
			http.MethodPatch,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusOK,
			[]byte(`Actor updated`),
		},
	}
	for _, test := range testTable {
		if test.name == "wrong method" {
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(MethodNotAllowed)).Return(0)
		} else if test.name == "wrong token" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(InvalidToken)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(InvalidToken)).Return(0)
		} else if test.name == "wrong json" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), mock.AnythingOfType("*json.UnmarshalTypeError")).Return(0)
		} else if test.name == "service error" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("UpdateActor", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("models.Actor")).Return(errors.New("service error")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("service error")).Return(0)
		} else if test.name == "positive" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("UpdateActor", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("models.Actor")).Return(nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string")).Return(0)
		}
		req, err := http.NewRequest(test.httpMethod, "/actor/update/"+test.id, bytes.NewBuffer(test.requestBody))
		req.Header.Set("Authorization", "Test-token")
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, req)
		responseBody := r.Body.Bytes()

		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, r.Code)
		assert.Equal(t, test.expectedResponseBody, responseBody)
	}
}
