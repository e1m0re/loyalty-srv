package handler

import (
	"encoding/json"
	"net/http"

	"e1m0re/loyalty-srv/internal/models"
)

type writingOffRequest struct {
	Order models.OrderNum `json:"order"`
	Sum   int             `json:"sum"`
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

	userID := models.UserID(1)
	account, err := h.services.AccountsService.GetAccountByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = h.services.AccountsService.Withdraw(r.Context(), account.ID, requestData.Sum, requestData.Order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
