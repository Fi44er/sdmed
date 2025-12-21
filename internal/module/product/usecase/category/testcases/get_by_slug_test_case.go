package category_testcases

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	"github.com/Fi44er/sdmed/internal/module/product/usecase/category/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MockGetBySlug struct {
	Ctrl     *gomock.Controller
	Ctx      context.Context
	RepoMock *mock.MockICategoryRepository
	FileMock *mock.MockIFileUsecaseAdapter
	T        assert.TestingT
}

type GetBySlugTestCase struct {
	Name             string
	InputSlug        string
	SetupMocks       func(m *MockGetBySlug)
	ExpectedError    error
	ExpectedCategory *product_entity.Category
}

func GetGetBySlugTestCases() []GetBySlugTestCase {
	return []GetBySlugTestCase{
		{
			Name:      "successful_get_by_slug",
			InputSlug: "test-slug",
			SetupMocks: func(m *MockGetBySlug) {
				category := &product_entity.Category{ID: "123", Name: "Test"}
				m.RepoMock.EXPECT().GetBySlug(m.Ctx, "test-slug").Return(category, nil)
				m.FileMock.EXPECT().GetByOwner(m.Ctx, "123", "category").Return([]product_entity.File{}, nil)
			},
			ExpectedError:    nil,
			ExpectedCategory: &product_entity.Category{ID: "123", Name: "Test", Images: []product_entity.File{}},
		},
		{
			Name:      "not_found",
			InputSlug: "unknown",
			SetupMocks: func(m *MockGetBySlug) {
				m.RepoMock.EXPECT().GetBySlug(m.Ctx, "unknown").Return(nil, nil)
			},
			ExpectedError:    product_constant.ErrCategoryNotFound,
			ExpectedCategory: nil,
		},
	}
}
