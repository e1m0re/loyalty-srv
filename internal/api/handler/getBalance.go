package handler

import (
	"encoding/json"
	"net/http"

	"e1m0re/loyalty-srv/internal/models"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(models.UserID)
	accountInfo, err := h.services.AccountsService.GetAccountInfoByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(accountInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
