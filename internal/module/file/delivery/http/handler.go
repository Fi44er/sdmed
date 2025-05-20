package http

import (
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type IFileUsecase interface{}

type FileHandler struct {
	usecase IFileUsecase

	logger    *logger.Logger
	validator *validator.Validate
}
