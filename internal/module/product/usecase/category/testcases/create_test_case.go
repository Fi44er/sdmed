package category_testcases

import (
	"context"
	"errors"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	"github.com/Fi44er/sdmed/internal/module/product/usecase/category/mock"
	uow_mock "github.com/Fi44er/sdmed/pkg/postgres/uow/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MockCreate struct {
	Ctrl               *gomock.Controller
	Ctx                context.Context
	RepoMock           *mock.MockICategoryRepository
	FileMock           *mock.MockIFileUsecaseAdapter
	CharacteristicMock *mock.MockICharacteristicUsecase
	UowMock            *uow_mock.MockUow
	T                  assert.TestingT
}

type CreateTestCase struct {
	Name          string
	InputCategory *product_entity.Category
	SetupMocks    func(m *MockCreate)
	ExpectedError error
}

func GetCreateTestCases() []CreateTestCase {
	return []CreateTestCase{
		{
			Name: "successful_creation_with_images",
			InputCategory: &product_entity.Category{
				ID:   "test-category-123",
				Name: "Test Category",
				Images: []product_entity.File{
					{Name: "image1.jpg"},
					{Name: "image2.png"},
				},
				Characteristics: []product_entity.Characteristic{
					{
						Name:       "test characteristic",
						CategoryID: "test-category-123",
					},
				},
			},
			SetupMocks: func(m *MockCreate) {
				m.UowMock.EXPECT().
					Do(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.UowMock.EXPECT().
					GetRepository(m.Ctx, "category").
					Return(m.RepoMock, nil)

				m.RepoMock.EXPECT().
					GetByName(m.Ctx, "Test Category").
					Return(nil, nil)

				m.RepoMock.EXPECT().
					Create(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, category *product_entity.Category) error {
						assert.Equal(m.T, "test-category-123", category.ID)
						assert.Equal(m.T, "Test Category", category.Name)
						return nil
					})

				m.CharacteristicMock.EXPECT().
					CreateMany(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, characteristics []product_entity.Characteristic) error {
						if len(characteristics) > 0 {
							assert.Equal(m.T, "test-category-123", characteristics[0].CategoryID)
							assert.Equal(m.T, "test characteristic", characteristics[0].Name)
						}
						return nil
					})

				m.FileMock.EXPECT().
					MakeFilesPermanent(m.Ctx, []string{"image1.jpg", "image2.png"}, "test-category-123", "category").
					Return(nil)
			},
			ExpectedError: nil,
		},
		{
			Name: "successful_creation_without_images",
			InputCategory: &product_entity.Category{
				ID:     "test-category-123",
				Name:   "Test Category",
				Images: []product_entity.File{},
				Characteristics: []product_entity.Characteristic{
					{
						Name:       "test characteristic",
						CategoryID: "test-category-123",
					},
				},
			},
			SetupMocks: func(m *MockCreate) {
				m.UowMock.EXPECT().
					Do(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.UowMock.EXPECT().
					GetRepository(m.Ctx, "category").
					Return(m.RepoMock, nil)

				m.RepoMock.EXPECT().
					GetByName(m.Ctx, "Test Category")

				m.RepoMock.EXPECT().
					Create(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, category *product_entity.Category) error {
						assert.Equal(m.T, "test-category-123", category.ID)
						assert.Equal(m.T, "Test Category", category.Name)
						return nil
					})

				m.CharacteristicMock.EXPECT().
					CreateMany(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, characteristics []product_entity.Characteristic) error {
						if len(characteristics) > 0 {
							assert.Equal(m.T, "test-category-123", characteristics[0].CategoryID)
							assert.Equal(m.T, "test characteristic", characteristics[0].Name)
						}
						return nil
					})
			},
			ExpectedError: nil,
		},
		{
			Name: "category already exists",
			InputCategory: &product_entity.Category{
				ID:   "test-category-456",
				Name: "Existing Category",
			},
			SetupMocks: func(m *MockCreate) {
				m.UowMock.EXPECT().
					Do(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.UowMock.EXPECT().
					GetRepository(m.Ctx, "category").
					Return(m.RepoMock, nil)

				existingCategory := &product_entity.Category{Name: "Existing Category"}
				m.RepoMock.EXPECT().
					GetByName(m.Ctx, "Existing Category").
					Return(existingCategory, nil)
			},
			ExpectedError: product_constant.ErrCategoryAlreadyExists,
		},
		{
			Name: "failed_to_create_category",
			InputCategory: &product_entity.Category{
				ID:   "test-category-123",
				Name: "Test Category",
				Images: []product_entity.File{
					{Name: "Test Image 1"},
				},
			},
			SetupMocks: func(m *MockCreate) {
				m.UowMock.EXPECT().
					Do(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.UowMock.EXPECT().
					GetRepository(m.Ctx, "category").
					Return(m.RepoMock, nil)

				m.RepoMock.EXPECT().
					GetByName(m.Ctx, "Test Category").
					Return(nil, nil)

				m.RepoMock.EXPECT().
					Create(m.Ctx, gomock.Any()).
					Return(errors.New("database error"))
			},
			ExpectedError: errors.New("database error"),
		},
		{
			Name: "failed_to_save_files",
			InputCategory: &product_entity.Category{
				ID:   "test-category-123",
				Name: "Test Category",
				Images: []product_entity.File{
					{Name: "image1.jpg"},
					{Name: "image2.png"},
				},
			},
			SetupMocks: func(m *MockCreate) {
				m.UowMock.EXPECT().
					Do(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				m.UowMock.EXPECT().
					GetRepository(m.Ctx, "category").
					Return(m.RepoMock, nil)

				m.RepoMock.EXPECT().
					GetByName(m.Ctx, "Test Category").
					Return(nil, nil)

				m.RepoMock.EXPECT().
					Create(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, category *product_entity.Category) error {
						assert.Equal(m.T, "test-category-123", category.ID)
						assert.Equal(m.T, "Test Category", category.Name)
						return nil
					})

				m.CharacteristicMock.EXPECT().
					CreateMany(m.Ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, characteristics []product_entity.Characteristic) error {
						if len(characteristics) > 0 {
							assert.Equal(m.T, "test-category-123", characteristics[0].CategoryID)
							assert.Equal(m.T, "test characteristic", characteristics[0].Name)
						}
						return nil
					})

				m.FileMock.EXPECT().
					MakeFilesPermanent(m.Ctx, []string{"image1.jpg", "image2.png"}, "test-category-123", "category").
					Return(errors.New("failed to save files"))
			},
			ExpectedError: errors.New("failed to save files"),
		},
	}
}
