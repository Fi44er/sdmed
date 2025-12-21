package category_testcases

import (
	"context"
	"errors"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/internal/module/product/usecase/category/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MockGetAll struct {
	Ctrl     *gomock.Controller
	Ctx      context.Context
	RepoMock *mock.MockICategoryRepository
	FileMock *mock.MockIFileUsecaseAdapter
	T        assert.TestingT
}

type GetAllTestCase struct {
	Name                    string
	InputOffset, InputLimit int
	SetupMocks              func(m *MockGetAll)
	ExpectedError           error
	ExpectedCategories      []product_entity.Category
}

func GetGetAllTestCases() []GetAllTestCase {
	ownerID1 := "test-category-123"
	ownerID2 := "test-category-456"
	category1 := product_entity.Category{
		ID:   "test-category-123",
		Name: "Test Category 1",
	}

	category2 := product_entity.Category{
		ID:   "test-category-456",
		Name: "Test Category 2",
	}

	filesCategory1 := []product_entity.File{
		{ID: "file-1", Name: "image1.jpg", OwnerID: &ownerID1},
		{ID: "file-2", Name: "image2.png", OwnerID: &ownerID1},
	}

	filesCategory2 := []product_entity.File{
		{ID: "file-3", Name: "image3.jpg", OwnerID: &ownerID2},
		{ID: "file-4", Name: "image4.png", OwnerID: &ownerID2},
	}

	return []GetAllTestCase{
		{
			Name:        "successful_get_all",
			InputOffset: 0,
			InputLimit:  10,
			SetupMocks: func(m *MockGetAll) {
				categories := []product_entity.Category{category1, category2}

				// 1. Ожидаем получение списка категорий
				m.RepoMock.EXPECT().
					GetAll(m.Ctx, 0, 10).
					Return(categories, nil)

				// 2. Ожидаем получение файлов
				ownerIDs := []string{"test-category-123", "test-category-456"}
				filesByOwner := map[string][]product_entity.File{
					"test-category-123": filesCategory1,
					"test-category-456": filesCategory2,
				}
				m.FileMock.EXPECT().
					GetByOwners(m.Ctx, ownerIDs, "category").
					Return(filesByOwner, nil)

				// 3. ИСПРАВЛЕНО: Добавлено ожидание Count, так как usecase вызывает его в конце
				m.RepoMock.EXPECT().
					Count(m.Ctx).
					Return(int64(2), nil)
			},
			ExpectedError: nil,
			ExpectedCategories: []product_entity.Category{
				{
					ID:     "test-category-123",
					Name:   "Test Category 1",
					Images: filesCategory1,
				},
				{
					ID:     "test-category-456",
					Name:   "Test Category 2",
					Images: filesCategory2,
				},
			},
		},
		{
			Name:        "successful_get_empty_categories_list",
			InputOffset: 0,
			InputLimit:  10,
			SetupMocks: func(m *MockGetAll) {
				categories := []product_entity.Category{}
				m.RepoMock.EXPECT().
					GetAll(m.Ctx, 0, 10).
					Return(categories, nil)

				// Здесь Count НЕ ожидается, так как usecase делает return, если len == 0
			},
			ExpectedError:      nil,
			ExpectedCategories: []product_entity.Category{},
		},
		{
			Name:        "failed_to_get_all_from_repository",
			InputOffset: 0,
			InputLimit:  10,
			SetupMocks: func(m *MockGetAll) {
				m.RepoMock.EXPECT().
					GetAll(m.Ctx, 0, 10).
					Return(nil, errors.New("failed to get categories"))

				// Здесь Count НЕ ожидается, так как usecase возвращает ошибку сразу
			},
			ExpectedError:      errors.New("failed to get categories"),
			ExpectedCategories: nil,
		},
		{
			Name:        "categories_loaded_but_files_failed",
			InputOffset: 0,
			InputLimit:  10,
			SetupMocks: func(m *MockGetAll) {
				categories := []product_entity.Category{category1, category2}
				m.RepoMock.EXPECT().
					GetAll(m.Ctx, 0, 10).
					Return(categories, nil)

				ownerIDs := []string{ownerID1, ownerID2}
				m.FileMock.EXPECT().
					GetByOwners(m.Ctx, ownerIDs, "category").
					Return(nil, errors.New("failed to get files"))

				// ИСПРАВЛЕНО: Даже если файлы упали, usecase логирует Warn и продолжает путь до Count
				m.RepoMock.EXPECT().
					Count(m.Ctx).
					Return(int64(2), nil)
			},
			ExpectedError: nil, // Так как usecase поглощает ошибку файлов через Warn
			ExpectedCategories: []product_entity.Category{
				{
					ID:     "test-category-123",
					Name:   "Test Category 1",
					Images: []product_entity.File(nil),
				},
				{
					ID:     "test-category-456",
					Name:   "Test Category 2",
					Images: []product_entity.File(nil),
				},
			},
		},
	}
}
