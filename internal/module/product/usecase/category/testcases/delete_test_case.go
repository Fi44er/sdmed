package category_testcases

import (
	"context"
	"errors"

	"github.com/Fi44er/sdmed/internal/module/product/usecase/category/mock"
	uow_mock "github.com/Fi44er/sdmed/pkg/postgres/uow/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MockDelete struct {
	Ctrl     *gomock.Controller
	Ctx      context.Context
	RepoMock *mock.MockICategoryRepository
	FileMock *mock.MockIFileUsecaseAdapter
	UowMock  *uow_mock.MockUow
	T        assert.TestingT
}

type DeleteTestCase struct {
	Name          string
	InputID       string
	SetupMocks    func(m *MockDelete)
	ExpectedError error
}

func GetDeleteTestCases() []DeleteTestCase {
	return []DeleteTestCase{
		{
			Name:    "successful_deletion",
			InputID: "cat-123",
			SetupMocks: func(m *MockDelete) {
				m.UowMock.EXPECT().Do(m.Ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
				m.UowMock.EXPECT().GetRepository(m.Ctx, "category").Return(m.RepoMock, nil)
				m.RepoMock.EXPECT().Delete(m.Ctx, "cat-123").Return(nil)
				m.FileMock.EXPECT().DeleteByOwner(m.Ctx, "cat-123", "category").Return(nil)
			},
			ExpectedError: nil,
		},
		{
			Name:    "failed_deletion_repo_error",
			InputID: "cat-123",
			SetupMocks: func(m *MockDelete) {
				m.UowMock.EXPECT().Do(m.Ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
				m.UowMock.EXPECT().GetRepository(m.Ctx, "category").Return(m.RepoMock, nil)
				m.RepoMock.EXPECT().Delete(m.Ctx, "cat-123").Return(errors.New("delete error"))
			},
			ExpectedError: errors.New("delete error"),
		},
	}
}
