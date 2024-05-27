package handler

import (
	"e1m0re/loyalty-srv/internal/apperrors"
	"encoding/json"
	"errors"
	"net/http"

	"e1m0re/loyalty-srv/internal/models"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	userInfo := &models.UserInfo{}
	err := json.NewDecoder(r.Body).Decode(userInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userInfo.Username == "" || userInfo.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.services.UsersService.SignIn(r.Context(), userInfo)
	if err != nil {
		if errors.Is(err, apperrors.ErrEntityNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(token) > 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}
