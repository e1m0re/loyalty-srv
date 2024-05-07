package handler

import (
	"errors"
	"net/http"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
)

func (handler handler) AddOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data []byte
	size, err := r.Body.Read(data)
	if err != nil || size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ordNum := models.OrderNum(data)
	_, isNewOrder, err := handler.OrderService.LoadOrder(r.Context(), ordNum)
	if err != nil {
		w.Write([]byte(err.Error()))

		switch true {
		case errors.Is(err, apperrors.InvalidOrderNumberError):
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, apperrors.OtherUsersOrderError):
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	if isNewOrder {
		w.WriteHeader(http.StatusAccepted)
	}
}
