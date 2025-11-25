package category_testcases

import (
	"context"
	"errors"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/internal/module/product/usecase/category/mock"
	uow_mock "github.com/Fi44er/sdmed/pkg/postgres/uow/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ownerID string = "test-category-123"

type MockGetByID struct {
	Ctrl     *gomock.Controller
	Ctx      context.Context
	RepoMock *mock.MockICategoryRepository
	FileMock *mock.MockIFileUsecaseAdapter
	UowMock  *uow_mock.MockUow
	T        assert.TestingT
}

type GetTestCase struct {
	Name             string
	InputID          string
	SetupMocks       func(m *MockGetByID)
	ExpectedError    error
	ExpectedCategory *product_entity.Category
}

func GetGetByIDTestCases() []GetTestCase {
	return []GetTestCase{
		{
			Name:    "successful_get",
			InputID: "test-category-123",
			SetupMocks: func(m *MockGetByID) {
				expectedCategory := &product_entity.Category{
					ID:   "test-category-123",
					Name: "Test Category",
				}
				expectedFiles := []product_entity.File{
					{ID: "file-1", Name: "image1.jpg", OwnerID: &ownerID},
					{ID: "file-2", Name: "image2.png", OwnerID: &ownerID},
				}

				m.RepoMock.EXPECT().
					GetByID(m.Ctx, "test-category-123").
					Return(expectedCategory, nil)

				m.FileMock.EXPECT().
					GetByOwner(m.Ctx, "test-category-123", "category").
					Return(expectedFiles, nil)
			},
			ExpectedError: nil,
			ExpectedCategory: &product_entity.Category{
				ID:   "test-category-123",
				Name: "Test Category",
				Images: []product_entity.File{
					{ID: "file-1", Name: "image1.jpg", OwnerID: &ownerID},
					{ID: "file-2", Name: "image2.png", OwnerID: &ownerID},
				},
			},
		},
		{
			Name:    "category_not_found",
			InputID: "non-existent-id",
			SetupMocks: func(m *MockGetByID) {
				m.RepoMock.EXPECT().
					GetByID(m.Ctx, "non-existent-id").
					Return(nil, errors.New("category not found"))
			},
			ExpectedError:    errors.New("category not found"),
			ExpectedCategory: nil,
		},
		{
			Name:    "failed_to_get_files_for_categor",
			InputID: "test-category-123",
			SetupMocks: func(m *MockGetByID) {
				expectedCategory := &product_entity.Category{
					ID:   "test-category-123",
					Name: "Test Category",
				}

				m.RepoMock.EXPECT().
					GetByID(m.Ctx, "test-category-123").
					Return(expectedCategory, nil)

				m.FileMock.EXPECT().
					GetByOwner(m.Ctx, "test-category-123", "category").
					Return(nil, assert.AnError)
			},
			ExpectedError: nil,
			ExpectedCategory: &product_entity.Category{
				ID:     "test-category-123",
				Name:   "Test Category",
				Images: []product_entity.File(nil),
			},
		},
	}
}
