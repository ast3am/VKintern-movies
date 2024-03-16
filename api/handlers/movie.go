package handlers

import (
	"encoding/json"
	"github.com/ast3am/VKintern-movies/internal/models"
	"io"
	"net/http"
	"net/url"
	"path"
)

func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, Admin)
	if err != nil {
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	var movie models.Movie
	err = movie.UnmarshalJSON(body)
	if err != nil {
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	err = h.services.CreateMovie(r.Context(), movie)
	if err != nil {
		if err.Error() == UnprocessableEntity {
			http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("movie created"))
}

func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, Admin)
	if err != nil {
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	var movie models.Movie
	err = movie.UnmarshalJSON(body)
	if err != nil {
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	id := path.Base(r.URL.String())

	err = h.services.UpdateMovie(r.Context(), id, movie)
	if err != nil {
		if err.Error() == UnprocessableEntity {
			http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("movie updated"))
}

func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, Admin)
	if err != nil {
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	id := path.Base(r.URL.String())

	err = h.services.DeleteMovie(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("movie deleted"))
}

func (h *Handler) GetMoviesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, User)
	if err != nil {
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	parseURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryParams, err := url.ParseQuery(parseURL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movieArr, err := h.services.GetMovieList(r.Context(), queryParams.Get("sortby"), queryParams.Get("line"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(movieArr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	err := h.services.CheckToken(token, User)
	if err != nil {
		http.Error(w, InvalidToken, http.StatusUnauthorized)
		return
	}

	parseURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryParams, err := url.ParseQuery(parseURL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movieArr, err := h.services.GetMovie(r.Context(), queryParams.Get("actor"), queryParams.Get("movie"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(movieArr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
