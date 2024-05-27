package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	"e1m0re/loyalty-srv/internal/repository/mocks"
)

func Test_ordersService_ValidateNumber(t *testing.T) {
	type args struct {
		ctx      context.Context
		orderNum models.OrderNum
	}
	type want struct {
		ok     bool
		errMsg string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty order number",
			args: args{
				ctx:      context.Background(),
				orderNum: "",
			},
			want: want{
				ok:     false,
				errMsg: apperrors.EmptyOrderNumberError.Error(),
			},
		},
		{
			name: "order number contains not only numbers",
			args: args{
				ctx:      context.Background(),
				orderNum: "123a-123",
			},
			want: want{
				ok:     false,
				errMsg: "strconv.Atoi: parsing \"a\": invalid syntax",
			},
		},
		{
			name: "valid order number",
			args: args{
				ctx:      context.Background(),
				orderNum: "12345678904",
			},
			want: want{
				ok: false,
			},
		},
		{
			name: "valid order number",
			args: args{
				ctx:      context.Background(),
				orderNum: "12345678903",
			},
			want: want{
				ok: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os := ordersService{
				orderRepository: mocks.NewOrderRepository(t),
			}
			gotOk, gotErr := os.ValidateNumber(test.args.ctx, test.args.orderNum)
			require.Equal(t, test.want.ok, gotOk)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}
