package service

import (
	"context"
	"fmt"
	"testing"

	"e1m0re/loyalty-srv/internal/apperrors"
	"e1m0re/loyalty-srv/internal/models"
	repositoriesMocks "e1m0re/loyalty-srv/internal/repository/mocks"
	servicesMocks "e1m0re/loyalty-srv/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_usersService_CreateUser(t *testing.T) {
	userInfo := models.UserInfo{
		Username: "user",
		Password: "password",
	}
	hashPassword := "$2a$10$LHhS8a.lcGHmhrgh6HiR0uLpmjG69OqrQ7YtTyPjapk1hlCF7hEhW"
	user := models.User{
		ID:       1,
		Username: userInfo.Username,
		Password: hashPassword,
	}

	type args struct {
		ctx      context.Context
		userInfo models.UserInfo
	}
	type want struct {
		user   *models.User
		errMsg string
	}
	tests := []struct {
		name           string
		prepareService func() *usersService
		args           args
		want           want
	}{
		{
			name: "GetPasswordHash failed",
			prepareService: func() *usersService {
				mockSecurityService := servicesMocks.NewSecurityService(t)
				mockSecurityService.
					On("GetPasswordHash", userInfo.Password).
					Return("", fmt.Errorf("some bcrypt error"))

				return &usersService{
					securityService: mockSecurityService,
				}
			},
			args: args{
				ctx:      context.Background(),
				userInfo: userInfo,
			},
			want: want{
				user:   nil,
				errMsg: "some bcrypt error",
			},
		},
		{
			name: "username is busy",
			prepareService: func() *usersService {
				mockSecurityService := servicesMocks.NewSecurityService(t)
				mockSecurityService.
					On("GetPasswordHash", userInfo.Password).
					Return(userInfo.Password, nil)

				mockUserRepo := repositoriesMocks.NewUserRepository(t)
				mockUserRepo.
					On("CreateUser", mock.Anything, userInfo).
					Return(nil, apperrors.BusyLoginError)

				return &usersService{
					userRepository:  mockUserRepo,
					securityService: mockSecurityService,
				}
			},
			args: args{
				ctx:      context.Background(),
				userInfo: userInfo,
			},
			want: want{
				user:   nil,
				errMsg: apperrors.BusyLoginError.Error(),
			},
		},
		{
			name: "CreateUser failed",
			prepareService: func() *usersService {
				mockSecurityService := servicesMocks.NewSecurityService(t)
				mockSecurityService.
					On("GetPasswordHash", userInfo.Password).
					Return(userInfo.Password, nil)

				mockUserRepo := repositoriesMocks.NewUserRepository(t)
				mockUserRepo.
					On("CreateUser", mock.Anything, userInfo).
					Return(nil, fmt.Errorf("some repos error"))

				return &usersService{
					userRepository:  mockUserRepo,
					securityService: mockSecurityService,
				}
			},
			args: args{
				ctx:      context.Background(),
				userInfo: userInfo,
			},
			want: want{
				user:   nil,
				errMsg: "some repos error",
			},
		},
		{
			name: "Successfully case",
			prepareService: func() *usersService {
				mockSecurityService := servicesMocks.NewSecurityService(t)
				mockSecurityService.
					On("GetPasswordHash", userInfo.Password).
					Return(hashPassword, nil)

				mockUserRepo := repositoriesMocks.NewUserRepository(t)
				mockUserRepo.
					On("CreateUser", mock.Anything, models.UserInfo{
						Username: userInfo.Username,
						Password: hashPassword,
					}).
					Return(&user, nil)

				return &usersService{
					userRepository:  mockUserRepo,
					securityService: mockSecurityService,
				}
			},
			args: args{
				ctx:      context.Background(),
				userInfo: userInfo,
			},
			want: want{
				user:   &user,
				errMsg: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			us := test.prepareService()
			gotUser, gotErr := us.CreateUser(test.args.ctx, &test.args.userInfo)
			require.Equal(t, &test.want.user, &gotUser)
			if len(test.want.errMsg) > 0 {
				require.EqualError(t, gotErr, test.want.errMsg)
			}
		})
	}
}
