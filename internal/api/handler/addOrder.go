package handler

import (
	"errors"
	"io"
	"net/http"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
)

func (h *Handler) AddOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ordNum := models.OrderNum(data)
	_, isNewOrder, err := h.services.Orders.NewOrder(r.Context(), ordNum)
	if err != nil {
		switch true {
		case errors.Is(err, apperrors.InvalidOrderNumberError):
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, apperrors.OtherUsersOrderError):
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	}

	if isNewOrder {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	w.WriteHeader(http.StatusOK)
}
