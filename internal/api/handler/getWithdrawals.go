package handler

import "net/http"

func (handler handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HandleGetWithdrawals"))
}
