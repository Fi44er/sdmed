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
	hash, _ := auth_utils.HashString(email)
	return hash
}

var testUser = auth_entity.User{
	PhoneNumber: "12345678901",
	Email:       "user@example.com",
	Password:    "hashedpassword",
	FIO:         "Test User",
}

var VerifyCodeTests = []struct {
	Name        string
	Input       *auth_entity.Code
	ExpectedErr error
	SetupMocks  func(ctx context.Context, m *MockVerifyCodeDeps)
}{
	{
		Name: "Success",
		Input: &auth_entity.Code{
			Code:  "123456",
			Email: testUser.Email,
			Type:  auth_entity.CodeTypeVerify,
		},
		ExpectedErr: nil,
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)

			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest any) error {
						*(dest.(*string)) = "123456"
						return nil
					}),
				m.Cache.EXPECT().Del(ctx, "verification_codes_"+hash).Return(nil),
				m.Cache.EXPECT().
					Get(ctx, "temp_user_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest any) error {
						*(dest.(*auth_entity.User)) = testUser
						return nil
					}),
				m.User.EXPECT().Create(ctx, gomock.Any()).Return(nil),
				m.Cache.EXPECT().Del(ctx, "temp_user_"+hash).Return(nil),
			)
		},
	},
	{
		Name:        "GetCodeError",
		Input:       &auth_entity.Code{Code: "123456", Email: testUser.Email, Type: auth_entity.CodeTypeVerify},
		ExpectedErr: auth_constant.ErrInternalServerError,
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			m.Cache.EXPECT().
				Get(ctx, "verification_codes_"+hash, gomock.Any()).
				Return(errors.New("redis error"))
		},
	},
	{
		Name:        "IncorrectVerificationCode",
		Input:       &auth_entity.Code{Code: "wrongcode", Email: testUser.Email, Type: auth_entity.CodeTypeVerify},
		ExpectedErr: auth_constant.ErrInternalServerError,
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			m.Cache.EXPECT().
				Get(ctx, "verification_codes_"+hash, gomock.Any()).
				DoAndReturn(func(_ context.Context, _ string, dest any) error {
					*(dest.(*string)) = "realcode"
					return auth_constant.ErrInternalServerError
				})
		},
	},
	{
		Name:        "DeleteCodeError",
		Input:       &auth_entity.Code{Code: "123456", Email: testUser.Email, Type: auth_entity.CodeTypeVerify},
		ExpectedErr: errors.New("delete code error"),
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest any) error {
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
		Input:       &auth_entity.Code{Code: "123456", Email: testUser.Email, Type: auth_entity.CodeTypeVerify},
		ExpectedErr: errors.New("temp user not found"),
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest any) error {
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
		Input:       &auth_entity.Code{Code: "123456", Email: testUser.Email, Type: auth_entity.CodeTypeVerify},
		ExpectedErr: errors.New("create user failed"),
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest any) error {
						*(dest.(*string)) = "123456"
						return nil
					}),
				m.Cache.EXPECT().Del(ctx, "verification_codes_"+hash).Return(nil),
				m.Cache.EXPECT().
					Get(ctx, "temp_user_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest any) error {
						*(dest.(*auth_entity.User)) = testUser
						return nil
					}),
				m.User.EXPECT().Create(ctx, gomock.Any()).
					Return(errors.New("create user failed")),
			)
		},
	},
	{
		Name:        "DeleteTempUserError",
		Input:       &auth_entity.Code{Code: "123456", Email: testUser.Email, Type: auth_entity.CodeTypeVerify},
		ExpectedErr: errors.New("cache delete failed"),
		SetupMocks: func(ctx context.Context, m *MockVerifyCodeDeps) {
			hash := getHashedEmail(testUser.Email)
			gomock.InOrder(
				m.Cache.EXPECT().
					Get(ctx, "verification_codes_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest any) error {
						*(dest.(*string)) = "123456"
						return nil
					}),
				m.Cache.EXPECT().Del(ctx, "verification_codes_"+hash).Return(nil),
				m.Cache.EXPECT().
					Get(ctx, "temp_user_"+hash, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, dest any) error {
						*(dest.(*auth_entity.User)) = testUser
						return nil
					}),
				m.User.EXPECT().Create(ctx, gomock.Any()).Return(nil),
				m.Cache.EXPECT().Del(ctx, "temp_user_"+hash).Return(errors.New("cache delete failed")),
			)
		},
	},
}
