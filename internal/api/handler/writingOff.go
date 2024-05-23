package handler

import (
	"encoding/json"
	"io"
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

	requestDataRaw, err := io.ReadAll(r.Body)
	if err != nil || len(requestDataRaw) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestData := writingOffRequest{}
	err = json.Unmarshal(requestDataRaw, &requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := h.services.Orders.ValidateNumber(r.Context(), requestData.Order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userId := models.UserId(1)
	account, err := h.services.Accounts.GetAccountByUserId(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = h.services.Accounts.Withdraw(r.Context(), account.ID, requestData.Sum, requestData.Order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
