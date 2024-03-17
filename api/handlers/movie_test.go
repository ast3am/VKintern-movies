package handlers

import (
	"bytes"
	"errors"
	"github.com/ast3am/VKintern-movies/api/handlers/mocks"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func GetDate(date string) time.Time {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, date)
	return t
}

func TestHandler_CreateMovie(t *testing.T) {
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
			[]byte(`{"name": "Forrest Gump", "description": "Description 1", "release_date": "1994-07-06", "rating": 8.5, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
			http.MethodGet,
			http.StatusMethodNotAllowed,
			[]byte(`invalid request method` + "\n"),
		}, {
			"wrong token",
			[]byte(`{"name": "Forrest Gump", "description": "Description 1", "release_date": "1994-07-06", "rating": 8.5, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
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
			[]byte(`{"name": "", "description": "Description 1", "release_date": "1994-07-06", "rating": 8.5, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
			http.MethodPost,
			http.StatusUnprocessableEntity,
			[]byte(`no valid data` + "\n"),
		}, {
			"another error service",
			[]byte(`{"name": "Forrest Gump", "description": "Description 1", "release_date": "1994-07-06", "rating": 8.5, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
			http.MethodPost,
			http.StatusBadRequest,
			[]byte(`service error` + "\n"),
		}, {
			"positive",
			[]byte(`{"name": "Forrest Gump", "description": "Description 1", "release_date": "1994-07-06", "rating": 8.5, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
			http.MethodPost,
			http.StatusOK,
			[]byte(`Movie created`),
		},
	}
	for _, test := range testTable {
		if test.name == "wrong method" {
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("invalid request method")).Return(0)
		} else if test.name == "wrong token" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(InvalidToken)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(InvalidToken)).Return(0)
		} else if test.name == "wrong json" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), mock.AnythingOfType("*json.UnmarshalTypeError")).Return(0)
		} else if test.name == "no valid data" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("CreateMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("models.Movie")).Return(errors.New(UnprocessableEntity)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(UnprocessableEntity)).Return(0)
		} else if test.name == "another error service" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("CreateMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("models.Movie")).Return(errors.New("service error")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("service error")).Return(0)
		} else if test.name == "positive" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("CreateMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("models.Movie")).Return(nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string")).Return(0)
		}
		req, err := http.NewRequest(test.httpMethod, "/movie/create", bytes.NewBuffer(test.requestBody))
		req.Header.Set("Authorization", "Test-token")
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, req)
		responseBody := r.Body.Bytes()

		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, r.Code)
		assert.Equal(t, test.expectedResponseBody, responseBody)
	}
}

func TestHandler_GetMoviesList(t *testing.T) {
	mux := http.NewServeMux()
	serv := mocks.NewService(t)
	log := mocks.NewLogger(t)
	h := NewHandler(serv, log)
	h.RegisterHandlers(mux)

	testTable := []struct {
		name                 string
		httpMethod           string
		urlValues            string
		expectedStatusCode   int
		expectedResponseBody []byte
	}{
		{
			"wrong method",
			http.MethodPost,
			"sortby=name&line=asc",
			http.StatusMethodNotAllowed,
			[]byte(`invalid request method` + "\n"),
		}, {
			"wrong token",
			http.MethodGet,
			"sortby=name&line=asc",
			http.StatusUnauthorized,
			[]byte(`invalid token` + "\n"),
		}, {
			"wrong service",
			http.MethodGet,
			"sortby=name&line=asc",
			http.StatusBadRequest,
			[]byte(`service error` + "\n"),
		}, {
			"positive",
			http.MethodGet,
			"sortby=name&line=asc",
			http.StatusOK,
			[]byte(`[{"name":"hatiko","description":"the movie about sad dog","release_date":"1992-04-01T00:00:00Z","rating":9.1,"actor_list":["suize"]}]`),
		},
	}
	for _, test := range testTable {
		if test.name == "wrong method" {
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("invalid request method")).Return(0).Once()
		} else if test.name == "wrong token" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(InvalidToken)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(InvalidToken)).Return(0).Once()
		} else if test.name == "wrong service" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
			serv.On("GetMovieList", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("service error")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("service error")).Return(0).Once()
		} else if test.name == "positive" {
			result := []*models.Movie{{Name: "hatiko", Description: "the movie about sad dog", ReleaseDate: GetDate("1992-04-01"), Rating: 9.1, ActorList: []string{"suize"}}}
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
			serv.On("GetMovieList", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(result, nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string")).Return(0).Once()
		}
		req, err := http.NewRequest(test.httpMethod, "/movie/get-list?"+test.urlValues, nil)
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

func TestHandler_GetMovies(t *testing.T) {
	mux := http.NewServeMux()
	serv := mocks.NewService(t)
	log := mocks.NewLogger(t)
	h := NewHandler(serv, log)
	h.RegisterHandlers(mux)

	testTable := []struct {
		name                 string
		httpMethod           string
		urlValues            string
		expectedStatusCode   int
		expectedResponseBody []byte
	}{
		{
			"wrong method",
			http.MethodPost,
			"actor=sui&movie=h",
			http.StatusMethodNotAllowed,
			[]byte(`invalid request method` + "\n"),
		}, {
			"wrong token",
			http.MethodGet,
			"actor=sui&movie=h",
			http.StatusUnauthorized,
			[]byte(`invalid token` + "\n"),
		}, {
			"wrong service",
			http.MethodGet,
			"actor=sui&movie=h",
			http.StatusBadRequest,
			[]byte(`service error` + "\n"),
		}, {
			"positive",
			http.MethodGet,
			"actor=sui&movie=h",
			http.StatusOK,
			[]byte(`[{"name":"hatiko","description":"the movie about sad dog","release_date":"1992-04-01T00:00:00Z","rating":9.1,"actor_list":["suize"]}]`),
		},
	}
	for _, test := range testTable {
		if test.name == "wrong method" {
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("invalid request method")).Return(0).Once()
		} else if test.name == "wrong token" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New(InvalidToken)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(InvalidToken)).Return(0).Once()
		} else if test.name == "wrong service" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
			serv.On("GetMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("service error")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("service error")).Return(0).Once()
		} else if test.name == "positive" {
			result := []*models.Movie{{Name: "hatiko", Description: "the movie about sad dog", ReleaseDate: GetDate("1992-04-01"), Rating: 9.1, ActorList: []string{"suize"}}}
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
			serv.On("GetMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(result, nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string")).Return(0).Once()
		}
		req, err := http.NewRequest(test.httpMethod, "/movie/get-movie?"+test.urlValues, nil)
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

func TestHandler_DeleteMovie(t *testing.T) {
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
			[]byte(`Movie deleted`),
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
			serv.On("DeleteMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).Return(errors.New("service problem")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("service problem")).Return(0).Once()
		} else if test.name == "positive" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
			serv.On("DeleteMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).Return(nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string")).Return(0).Once()
		}
		req, err := http.NewRequest(test.httpMethod, "/movie/delete/"+test.id, nil)
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
func TestHandler_UpdateMovie(t *testing.T) {
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
			[]byte(`{"name": "Borrest Gump", "description": "Description 1", "release_date": "1994-07-06", "rating": 8.5, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
			http.MethodGet,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusMethodNotAllowed,
			[]byte(`invalid request method` + "\n"),
		}, {
			"wrong token",
			[]byte(`{"name": "Borrest Gump", "description": "Description 1", "release_date": "1994-07-06", "rating": 8.5, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
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
			"no valid data",
			[]byte(`{"name": "Borrest Gump", "description": "Description 1", "release_date": "1994-07-06", "rating": 12, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
			http.MethodPatch,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusUnprocessableEntity,
			[]byte(`no valid data` + "\n"),
		}, {
			"service error",
			[]byte(`{"name": "Borrest Gump", "description": "Description 1", "release_date": "1994-07-06", "rating": 8.5, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
			http.MethodPatch,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusBadRequest,
			[]byte(`service error` + "\n"),
		}, {
			"positive",
			[]byte(`{"name": "Borrest Gump", "description": "Description 1", "release_date": "1994-07-06", "rating": 8.5, "actor_list": ["Tom Hanks", "Brad Pitt"]}`),
			http.MethodPatch,
			"40d882f7-b027-4a07-85da-76e0f7d9b6e3",
			http.StatusOK,
			[]byte(`Movie updated`),
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
			serv.On("UpdateMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("models.Movie")).Return(errors.New(UnprocessableEntity)).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New(UnprocessableEntity)).Return(0)
		} else if test.name == "service error" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("UpdateMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("models.Movie")).Return(errors.New("service error")).Once()
			log.On("HandlerErrorLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string"), errors.New("service error")).Return(0)
		} else if test.name == "positive" {
			serv.On("CheckToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			serv.On("UpdateMovie", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("models.Movie")).Return(nil).Once()
			log.On("HandlerLog", mock.AnythingOfType("*http.Request"), test.expectedStatusCode, mock.AnythingOfType("string")).Return(0)
		}
		req, err := http.NewRequest(test.httpMethod, "/movie/update/"+test.id, bytes.NewBuffer(test.requestBody))
		req.Header.Set("Authorization", "Test-token")
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, req)
		responseBody := r.Body.Bytes()

		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, r.Code)
		assert.Equal(t, test.expectedResponseBody, responseBody)
	}
}
