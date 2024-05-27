package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userInfo := &models.UserInfo{}
	err := json.NewDecoder(r.Body).Decode(userInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if userInfo.Username == "" || userInfo.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.services.UsersService.CreateUser(r.Context(), userInfo)
	if err != nil {
		if errors.Is(err, apperrors.ErrBusyLogin) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.services.UsersService.SignIn(r.Context(), userInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
