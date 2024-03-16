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
	w.Write([]byte("movie deleted"))
	h.log.HandlerLog(r, http.StatusOK, "Movie deleted")
}

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
