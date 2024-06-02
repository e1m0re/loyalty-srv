package handler

import (
	"errors"
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

	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("AddOrder", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderInfo := models.OrderInfo{
		UserID:   r.Context().Value(models.CKUserID).(models.UserID),
		OrderNum: models.OrderNum(data),
	}

	_, err = h.services.OrdersService.NewOrder(r.Context(), orderInfo)
	if err != nil {
		switch true {
		case errors.Is(err, apperrors.ErrInvalidOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, apperrors.ErrOrderWasLoadedByAnotherUser):
			w.WriteHeader(http.StatusConflict)
		case errors.Is(err, apperrors.ErrOrderWasLoaded):
			w.WriteHeader(http.StatusOK)
		default:
			slog.Error("AddOrder", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
