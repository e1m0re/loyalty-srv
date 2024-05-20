package handler

import (
	"encoding/json"
	"net/http"

	"e1m0re/loyalty-srv/internal/models"
)

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userId := models.UserId(1)
	withdrawalsList, err := h.services.Accounts.GetWithdrawals(r.Context(), userId)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(withdrawalsList)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
	w.WriteHeader(http.StatusOK)
}
