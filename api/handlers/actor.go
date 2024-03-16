package handlers

import (
	"encoding/json"
	"github.com/ast3am/VKintern-movies/internal/models"
	"io"
	"net/http"
	"path"
)

func (h *Handler) CreateActor(w http.ResponseWriter, r *http.Request) {
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

	var actor models.Actor
	err = actor.UnmarshalJSON(body)
	if err != nil {
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	err = h.services.CreateActor(r.Context(), actor)
	if err.Error() == "empty name" {
		http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("actor created"))
}

func (h *Handler) GetActorsList(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.services.GetActorList(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
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

	var actor models.Actor
	err = actor.UnmarshalJSON(body)
	if err != nil {
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	id := path.Base(r.URL.String())

	err = h.services.UpdateActor(r.Context(), id, actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("actor updated"))
}

func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
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

	err = h.services.DeleteActor(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("actor deleted"))
}
