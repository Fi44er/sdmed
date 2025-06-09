package setup_mocks

import (
	"context"
	"errors"

	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/utils"
	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
	"github.com/golang/mock/gomock"
)

type MockVerifyCodeDeps struct {
	User  *mocks.MockIUserUsecase
	Cache *mocks.MockICache
}

func getHashedEmail(email string) string {
	hash, _ := utils.HashString(email)
	return hash
}

var testUser = entity.User{
	PhoneNumber: "12345678901",
	Email:       "user@example.com",
	Password:    "hashedpassword",
	FIO:         "Test User",
}

var VerifyCodeTests = []struct {
	Name        string
	Input       *entity.VerifyCode
	ExpectedErr error
	SetupMocks  func(ctx context.Context, m *MockVerifyCodeDeps)
}{
	{
		Name: "Success",
		Input: &entity.VerifyCode{
			Code:  "123456",
			Email: testUser.Email,
		},
		ExpectedErr: nil,
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)

			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
						*(dest.(*string)) = "123456"
						return nil
					}),
				m.Cache.EXPECT().Del(ctx, "verification_codes_"+hash).Return(nil),
				m.Cache.EXPECT().
					Get(ctx, "temp_user_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
						*(dest.(*entity.User)) = testUser
						return nil
					}),
				m.User.EXPECT().Create(ctx, gomock.Any()).Return(nil),
				m.Cache.EXPECT().Del(ctx, "temp_user_"+hash).Return(nil),
			)
		},
	},
	{
		Name:        "GetCodeError",
		Input:       &entity.VerifyCode{Code: "123456", Email: testUser.Email},
		ExpectedErr: constant.ErrInternalServerError,
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			m.Cache.EXPECT().
				Get(ctx, "verification_codes_"+hash, gomock.Any()).
				Return(errors.New("redis error"))
		},
	},
	{
		Name:        "IncorrectVerificationCode",
		Input:       &entity.VerifyCode{Code: "wrongcode", Email: testUser.Email},
		ExpectedErr: constant.ErrInternalServerError,
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			m.Cache.EXPECT().
				Get(ctx, "verification_codes_"+hash, gomock.Any()).
				DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
					*(dest.(*string)) = "realcode"
					return constant.ErrInternalServerError
				})
		},
	},
	{
		Name:        "DeleteCodeError",
		Input:       &entity.VerifyCode{Code: "123456", Email: testUser.Email},
		ExpectedErr: errors.New("delete code error"),
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
						*(dest.(*string)) = "123456"
						return nil
					}),
				m.Cache.EXPECT().
					Del(ctx, "verification_codes_"+hash).
					Return(errors.New("delete code error")),
			)
		},
	},
	{
		Name:        "GetTempUserError",
		Input:       &entity.VerifyCode{Code: "123456", Email: testUser.Email},
		ExpectedErr: errors.New("temp user not found"),
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
						*(dest.(*string)) = "123456"
						return nil
					}),
				m.Cache.EXPECT().Del(ctx, "verification_codes_"+hash).Return(nil),
				m.Cache.EXPECT().
					Get(ctx, "temp_user_"+hash, gomock.Any()).
					Return(errors.New("temp user not found")),
			)
		},
	},
	{
		Name:        "CreateUserError",
		Input:       &entity.VerifyCode{Code: "123456", Email: testUser.Email},
		ExpectedErr: errors.New("create user failed"),
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
						*(dest.(*string)) = "123456"
						return nil
					}),
				m.Cache.EXPECT().Del(ctx, "verification_codes_"+hash).Return(nil),
				m.Cache.EXPECT().
					Get(ctx, "temp_user_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
						*(dest.(*entity.User)) = testUser
						return nil
					}),
				m.User.EXPECT().Create(ctx, gomock.Any()).
					Return(errors.New("create user failed")),
			)
		},
	},
	{
		Name:        "DeleteTempUserError",
		Input:       &entity.VerifyCode{Code: "123456", Email: testUser.Email},
		ExpectedErr: errors.New("cache delete failed"),
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
						*(dest.(*string)) = "123456"
						return nil
					}),
				m.Cache.EXPECT().Del(ctx, "verification_codes_"+hash).Return(nil),
				m.Cache.EXPECT().
					Get(ctx, "temp_user_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
						*(dest.(*entity.User)) = testUser
						return nil
					}),
				m.User.EXPECT().Create(ctx, gomock.Any()).Return(nil),
				m.Cache.EXPECT().Del(ctx, "temp_user_"+hash).Return(errors.New("cache delete failed")),
			)
		},
	},
}
