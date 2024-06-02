package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"e1m0re/loyalty-srv/internal/models"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(models.CKUserID).(models.UserID)
	invoice, err := h.services.InvoicesService.GetInvoiceByUserID(r.Context(), userID)
	if err != nil {
		slog.Error("GetBalance", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	invoiceInfo, err := h.services.InvoicesService.GetInvoiceInfo(r.Context(), invoice)
	if err != nil {
		slog.Error("GetBalance", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(invoiceInfo)
	if err != nil {
		slog.Error("GetBalance", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
