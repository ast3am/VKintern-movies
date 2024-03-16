package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ast3am/VKintern-movies/internal/models"
	"io"
	"net/http"
	"path"
)

func (h *Handler) CreateActor(w http.ResponseWriter, r *http.Request) {
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

	var actor models.Actor
	err = actor.UnmarshalJSON(body)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, ParsingJSONError, err)
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	err = h.services.CreateActor(r.Context(), actor)
	if err != nil {
		if err.Error() == UnprocessableEntity {
			h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, "Create actor", err)
			http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
			return
		} else {
			h.log.HandlerErrorLog(r, http.StatusBadRequest, "Create actor", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Actor created"))
	h.log.HandlerLog(r, http.StatusOK, "Actor created")
}

func (h *Handler) GetActorsList(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.services.GetActorList(r.Context())
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Get actor list", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusInternalServerError, "Can't marshal result", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	h.log.HandlerLog(r, http.StatusOK, "Actor list created")
}

func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
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

	var actor models.Actor
	err = actor.UnmarshalJSON(body)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, ParsingJSONError, err)
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	id := path.Base(r.URL.String())

	err = h.services.UpdateActor(r.Context(), id, actor)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Update actor", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("actor updated"))
	h.log.HandlerLog(r, http.StatusOK, "Actor updated")
}

func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
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

	err = h.services.DeleteActor(r.Context(), id)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "Delete actor", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Actor deleted"))
	h.log.HandlerLog(r, http.StatusOK, "Actor deleted")
}
