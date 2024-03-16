package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type UserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Auth godoc
// @Summary Вход в систему
// @Description Вход в систему по логину и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param data body UserDTO true "Входные параметры"
// @Success 200 {object} string
// @Failure 400,405,422 {object} error
// @Router /auth [post]
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.log.HandlerErrorLog(r, http.StatusMethodNotAllowed, "", errors.New(MethodNotAllowed))
		http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, UnprocessableEntity, err)
		http.Error(w, UnprocessableEntity, http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	var user UserDTO
	err = json.Unmarshal(body, &user)
	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusUnprocessableEntity, "", errors.New(ParsingJSONError))
		http.Error(w, ParsingJSONError, http.StatusUnprocessableEntity)
		return
	}

	token, err := h.services.Auth(r.Context(), user.Email, user.Password)

	if err != nil {
		h.log.HandlerErrorLog(r, http.StatusBadRequest, "", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"token": token,
	}
	w.Write([]byte("authorization success"))
	json.NewEncoder(w).Encode(response)
	h.log.HandlerLog(r, http.StatusOK, "authorization")
}
