package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ast3am/VKintern-movies/internal/models"
	"io"
	"net/http"
	"net/url"
	"path"
)

// CreateMovie godoc
// @Summary Создание фильма
// @Description Создание фильма, предполагается что все поля не пустые
// @Tags movie
// @Accept json
// @Produce json
// @Param data body models.Movie true "Входные параметры"
// @Success 200 {object} string
// @Failure 400,401,405,422 {object} error
// @Router /movie/create [post]
// @Security ApiKeyAuth
func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.log.HandlerErrorLog(r, http.StatusMethodNotAllowed, "", errors.New(MethodNotAllowed))
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, Admin)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnauthorized, InvalidToken, err)
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, UnprocessableEntity, err)
		http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	var movie models.Movie
	err = movie.UnmarshalJSON(body)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, ParsingJSONError, err)
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	err = h.services.CreateMovie(r.Context(), movie)
	if err != nil {
		if err.Error() == UnprocessableEntity {
			h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, "Create movie", err)
			http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
			return
		} else {
			h.log.HandlerErrorLog(r, http.StatusBadRequest, "Create movie", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Movie created"))
	h.log.HandlerLog(r, http.StatusOK, "Movie created")
}

// UpdateMovie godoc
// @Summary Изменение информации о фильме
// @Description Информация может быть изменена как частично, так и полностью
// @Tags movie
// @Accept json
// @Produce json
// @Param data body models.Movie true "Входные параметры"
// @Success 200 {object} string
// @Failure 400,401,405,422 {object} error
// @Router /movie/update/ [patch]
// @Security ApiKeyAuth
func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		h.log.HandlerErrorLog(r, http.StatusMethodNotAllowed, "", errors.New(MethodNotAllowed))
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, Admin)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnauthorized, InvalidToken, err)
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, UnprocessableEntity, err)
		http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	var movie models.Movie
	err = movie.UnmarshalJSON(body)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, ParsingJSONError, err)
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	id := path.Base(r.URL.String())

	err = h.services.UpdateMovie(r.Context(), id, movie)
	if err != nil {
		if err.Error() == UnprocessableEntity {
			h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, "Update movie", err)
			http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
			return
		} else {
			h.log.HandlerErrorLog(r, http.StatusBadRequest, "Update movie", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Movie updated"))
	h.log.HandlerLog(r, http.StatusOK, "Movie updated")
}

// DeleteMovie godoc
// @Summary Удаление информации о фильме
// @Description Полное удаление информации по UUID
// @Tags movie
// @Accept json
// @Produce json
// @Param uuid query string true "UUID фильма"
// @Success 200 {object} string
// @Failure 400,401,405 {object} error
// @Router /movie/delete [delete]
// @Security ApiKeyAuth
func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.log.HandlerErrorLog(r, http.StatusMethodNotAllowed, "", errors.New(MethodNotAllowed))
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, Admin)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnauthorized, InvalidToken, err)
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	id := path.Base(r.URL.String())

	err = h.services.DeleteMovie(r.Context(), id)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Delete actor", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Movie deleted"))
	h.log.HandlerLog(r, http.StatusOK, "Movie deleted")
}

// GetMoviesList godoc
// @Summary Получение списка фильмов
// @Description Получение списка с возможностью сортировки, параметры задаются в URL
// @Tags movie
// @Accept json
// @Produce json
// @Param sortby query string false "Указания поля для сортировки, по умолчанию rating)"
// @Param line query string false "Указание типа сортировки, по умолчанию desc"
// @Success 200 {object} models.Movie
// @Failure 400,401,405 {object} error
// @Failure 500 {object} error
// @Router /movie/get-list [get]
// @Security ApiKeyAuth
func (h *Handler) GetMoviesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.log.HandlerErrorLog(r, http.StatusMethodNotAllowed, "", errors.New(MethodNotAllowed))
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, User)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnauthorized, InvalidToken, err)
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	parseURL, err := url.Parse(r.URL.String())
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Get movie list URL parse error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryParams, err := url.ParseQuery(parseURL.RawQuery)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Get movie list URL parse error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movieArr, err := h.services.GetMovieList(r.Context(), queryParams.Get("sortby"), queryParams.Get("line"))
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Get movie list", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(movieArr)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusInternalServerError, "Can't marshal result", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	h.log.HandlerLog(r, http.StatusOK, "Movie list created")
}

// GetMovie godoc
// @Summary Получение списка фильмов
// @Description Получение списка фильмов с поиском по фрагменту названия и фрагменту имени актера
// @Tags movie
// @Accept json
// @Produce json
// @Param actor query string false "Указание актера"
// @Param movie query string false "Указание названия фильма"
// @Success 200 {object} models.Movie
// @Failure 400,401,405 {object} error
// @Failure 500 {object} error
// @Router /movie/get-movie [get]
// @Security ApiKeyAuth
func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.log.HandlerErrorLog(r, http.StatusMethodNotAllowed, "", errors.New(MethodNotAllowed))
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, User)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnauthorized, InvalidToken, err)
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	parseURL, err := url.Parse(r.URL.String())
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Get movie URL parse error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryParams, err := url.ParseQuery(parseURL.RawQuery)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Get movie URL parse error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movieArr, err := h.services.GetMovie(r.Context(), queryParams.Get("actor"), queryParams.Get("movie"))
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Get movie list", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(movieArr)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusInternalServerError, "Can't marshal result", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	h.log.HandlerLog(r, http.StatusOK, "Movie list by name created")
}
