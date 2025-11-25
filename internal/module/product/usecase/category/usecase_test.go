package category_usecase_test

import (
	"context"
	"testing"

	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
	category_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/category"
	"github.com/Fi44er/sdmed/internal/module/product/usecase/category/mock"
	category_testcases "github.com/Fi44er/sdmed/internal/module/product/usecase/category/testcases"
	"github.com/Fi44er/sdmed/pkg/logger"
	uow_mock "github.com/Fi44er/sdmed/pkg/postgres/uow/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CategoryUsecaseTestSuite struct {
	suite.Suite
	ctx      context.Context
	ctrl     *gomock.Controller
	usecase  category_usecase.ICategoryUsecase
	repoMock *mock.MockICategoryRepository
	fileMock *mock.MockIFileUsecaseAdapter
	uowMock  *uow_mock.MockUow
	logger   *logger.Logger
}

func (s *CategoryUsecaseTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.repoMock = mock.NewMockICategoryRepository(s.ctrl)
	s.fileMock = mock.NewMockIFileUsecaseAdapter(s.ctrl)
	s.uowMock = uow_mock.NewMockUow(s.ctrl)
	s.logger = logger.NewLogger()
	s.usecase = category_usecase.NewCategoryUsecase(s.logger, s.repoMock, s.fileMock, s.uowMock)
}

func TestCategoryUsecase(t *testing.T) {
	suite.Run(t, new(CategoryUsecaseTestSuite))
}

func (s *CategoryUsecaseTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *CategoryUsecaseTestSuite) TestCreate() {
	tests := category_testcases.GetCreateTestCases()

	for _, tc := range tests {
		s.T().Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := mock.NewMockICategoryRepository(ctrl)
			fileMock := mock.NewMockIFileUsecaseAdapter(ctrl)
			uowMock := uow_mock.NewMockUow(ctrl)

			mockStruct := &category_testcases.MockCreate{
				Ctrl:     ctrl,
				Ctx:      s.ctx,
				RepoMock: repoMock,
				FileMock: fileMock,
				UowMock:  uowMock,
				T:        t,
			}

			usecase := category_usecase.NewCategoryUsecase(s.logger, repoMock, fileMock, uowMock)

			tc.SetupMocks(mockStruct)

			err := usecase.Create(s.ctx, tc.InputCategory)

			if tc.ExpectedError != nil {
				assert.Error(t, err)
				if tc.ExpectedError != product_constant.ErrCategoryAlreadyExist {
					assert.Contains(t, err.Error(), tc.ExpectedError.Error())
				} else {
					assert.Equal(t, tc.ExpectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (s *CategoryUsecaseTestSuite) TestGetByID() {
	tests := category_testcases.GetGetByIDTestCases()

	for _, tc := range tests {
		s.T().Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := mock.NewMockICategoryRepository(ctrl)
			fileMock := mock.NewMockIFileUsecaseAdapter(ctrl)

			usecase := category_usecase.NewCategoryUsecase(s.logger, repoMock, fileMock, s.uowMock)

			mockStruct := &category_testcases.MockGetByID{
				Ctrl:     ctrl,
				Ctx:      s.ctx,
				RepoMock: repoMock,
				FileMock: fileMock,
				T:        t,
			}

			tc.SetupMocks(mockStruct)

			result, err := usecase.GetByID(s.ctx, tc.InputID)

			if tc.ExpectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.ExpectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.ExpectedCategory, result)
			}
		})
	}
}
