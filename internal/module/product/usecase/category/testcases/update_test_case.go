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

type MockUpdate struct {
	Ctrl               *gomock.Controller
	Ctx                context.Context
	RepoMock           *mock.MockICategoryRepository
	FileMock           *mock.MockIFileUsecaseAdapter
	CharacteristicMock *mock.MockICharacteristicUsecase
	UowMock            *uow_mock.MockUow
	T                  assert.TestingT
}

type UpdateTestCase struct {
	Name          string
	InputCategory *product_entity.Category
	SetupMocks    func(m *MockUpdate)
	ExpectedError error
}

func GetUpdateTestCases() []UpdateTestCase {
	categoryID := "cat-123"
	return []UpdateTestCase{
		{
			Name: "successful_full_update",
			InputCategory: &product_entity.Category{
				ID:   categoryID,
				Name: "New Name",
				Images: []product_entity.File{
					{Name: "new_image.jpg"}, // Это новое изображение (без ID)
				},
			},
			SetupMocks: func(m *MockUpdate) {
				m.UowMock.EXPECT().Do(m.Ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
				m.UowMock.EXPECT().GetRepository(m.Ctx, "category").Return(m.RepoMock, nil)

				// Возвращаем старую категорию (без изображений в структуре самого объекта)
				oldCategory := &product_entity.Category{ID: categoryID, Name: "Old Name"}
				m.RepoMock.EXPECT().GetByID(m.Ctx, categoryID).Return(oldCategory, nil)

				// Важно: ожидаем сохранение категории.
				// Если usecase обновляет объект, переданный в GetByID, gomock.Any() это поймает.
				m.RepoMock.EXPECT().Update(m.Ctx, gomock.Any()).Return(nil)

				// Имитируем, что в БД сейчас привязан файл "old_image.jpg"
				oldFiles := []product_entity.File{
					{ID: "old-f-1", Name: "old_image.jpg"},
				}
				m.FileMock.EXPECT().GetByOwner(m.Ctx, categoryID, "category").Return(oldFiles, nil)

				m.CharacteristicMock.EXPECT().CreateMany(m.Ctx, gomock.Any()).Return(nil)

				// Логика удаления старого файла (так как его нет в InputCategory)
				m.FileMock.EXPECT().DeleteByID(m.Ctx, "old-f-1").Return(nil)

				// ИСПРАВЛЕНИЕ:
				// Если тест все равно выдает [], значит в usecase.go:118
				// расчет новых имен файлов (new_image.jpg) не срабатывает.
				// Чтобы тест был гибким к реализации сравнения слайсов:
				m.FileMock.EXPECT().
					MakeFilesPermanent(m.Ctx, gomock.Eq([]string{"new_image.jpg"}), categoryID, "category").
					Return(nil)
			},
			ExpectedError: nil,
		},
		{
			Name:          "failed_repository_error",
			InputCategory: &product_entity.Category{ID: categoryID},
			SetupMocks: func(m *MockUpdate) {
				m.UowMock.EXPECT().Do(m.Ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
				m.UowMock.EXPECT().GetRepository(m.Ctx, "category").Return(m.RepoMock, nil)
				m.RepoMock.EXPECT().GetByID(m.Ctx, categoryID).Return(nil, errors.New("db error"))
			},
			ExpectedError: errors.New("db error"),
		},
	}
}
