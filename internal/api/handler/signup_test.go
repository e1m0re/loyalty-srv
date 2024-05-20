package handler

import (
	"net/http"
	"testing"

	"e1m0re/loyalty-srv/internal/service"
)

func Test_handler_SignUp(t *testing.T) {
	type fields struct {
		UserService    service.UsersService
		OrderService   service.OrdersService
		AccountService service.AccountsService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := handler{
				UserService:    tt.fields.UserService,
				OrderService:   tt.fields.OrderService,
				AccountService: tt.fields.AccountService,
			}
			handler.SignUp(tt.args.w, tt.args.r)
		})
	}
}
