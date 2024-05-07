package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
)

func (handler handler) SignUp(w http.ResponseWriter, r *http.Request) {
	userInfo := &models.UserInfo{}
	err := json.NewDecoder(r.Body).Decode(userInfo)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userInfo.Username == "" || userInfo.Password == "" {
		w.Write([]byte("invalid request"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = handler.UserService.SignUp(r.Context(), *userInfo)
	if err != nil {
		w.Write([]byte(err.Error()))
		if errors.Is(err, apperrors.BusyLoginError) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
