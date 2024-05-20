package handler

import (
	"encoding/json"
	"net/http"

	"e1m0re/loyalty-srv/internal/models"
)

func (handler handler) SignIn(w http.ResponseWriter, r *http.Request) {
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

	ok, err := handler.UserService.SignIn(r.Context(), *userInfo)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
