package handlers

import (
	"fmt"
	"github.com/ast3am/VKintern-movies/internal/models"
	"io"
	"net/http"
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
	if err.Error() == "No Valid Data" {
		http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("movie created"))
}

func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Movie Update")
}

func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Movie Delete")
}

func (h *Handler) GetMoviesList(w http.ResponseWriter, request *http.Request) {
	fmt.Println("Movie Get")
}

func (h *Handler) GetMovie(w http.ResponseWriter, request *http.Request) {
	fmt.Println("Movie Get one")
}
