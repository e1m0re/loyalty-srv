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

func (handler handler) WritingOff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requestDataRaw []byte
	size, err := r.Body.Read(requestDataRaw)
	if err != nil || size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestData := writingOffRequest{}
	err = json.Unmarshal(requestDataRaw, &requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := handler.OrderService.ValidateNumber(r.Context(), requestData.Order)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userId := models.UserId(1)
	account, err := handler.AccountService.GetAccountByUserId(r.Context(), userId)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = handler.AccountService.Withdraw(r.Context(), account.ID, requestData.Sum, requestData.Order)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
