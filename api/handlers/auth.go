package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

type UserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	var user UserDTO
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	token, err := h.services.Auth(r.Context(), user.Email, user.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"token": token,
	}
	w.Write([]byte("authorization success"))
	json.NewEncoder(w).Encode(response)
}
