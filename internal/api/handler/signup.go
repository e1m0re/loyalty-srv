package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userInfo := models.UserInfo{}
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userInfo.Username == "" || userInfo.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.services.UsersService.CreateUser(r.Context(), userInfo)
	if err != nil {
		if errors.Is(err, apperrors.ErrBusyLogin) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.services.AccountsService.CreateAccount(r.Context(), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := h.services.UsersService.SignIn(r.Context(), userInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}
