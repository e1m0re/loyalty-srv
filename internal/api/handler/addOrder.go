package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
)

func (h *Handler) AddOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(models.UserID)
	slog.Info("asdasd", slog.String("asd", fmt.Sprintf("%v", userID)))

	data, err := io.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ordNum := models.OrderNum(data)
	_, isNewOrder, err := h.services.OrdersService.NewOrder(r.Context(), ordNum)
	if err != nil {
		switch true {
		case errors.Is(err, apperrors.ErrInvalidOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, apperrors.ErrOtherUsersOrder):
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
