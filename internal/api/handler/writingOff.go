package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
)

type writingOffRequest struct {
	Order models.OrderNum `json:"order"`
	Sum   float64         `json:"sum"`
}

func (h *Handler) WritingOff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestData := &writingOffRequest{}
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := h.services.OrdersService.ValidateNumber(r.Context(), requestData.Order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userID := r.Context().Value(models.CKUserID).(models.UserID)
	account, err := h.services.AccountsService.GetAccountByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.services.AccountsService.UpdateBalance(r.Context(), *account, requestData.Sum, requestData.Order)
	if err != nil {
		if errors.Is(err, apperrors.ErrAccountHasNotEnoughFunds) {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
