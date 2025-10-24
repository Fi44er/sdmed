package http

import (
	"context"
	"io"

	"github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/internal/module/file/pkg/utils"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IFileUsecase interface {
	Upload(ctx context.Context, file *entity.File) error
	Get(ctx context.Context, name string) (*entity.File, error)
}

type FileHandler struct {
	usecase IFileUsecase

	logger    *logger.Logger
	validator *validator.Validate
}

func NewFileHandler(usecase IFileUsecase, logger *logger.Logger, validator *validator.Validate) *FileHandler {
	return &FileHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
	}
}

func (h *FileHandler) Get(ctx *fiber.Ctx) error {
	name := ctx.Params("name")

	file, err := h.usecase.Get(ctx.Context(), name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ctx.Set("Content-Disposition", "attachment; filename="+file.Name)
	ctx.Set("Content-Type", "application/octet-stream")

	contentType := utils.GetContentType(file.Name)
	ctx.Set("Content-Type", contentType)

	return ctx.Send(file.Data)
}

func (h *FileHandler) Upload(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file content",
		})
	}

	uploadDTO := &entity.File{
		OwnerID:   ctx.FormValue("owner_id"),
		OwnerType: ctx.FormValue("owner_type"),
		Data:      fileData,
	}

	if err := h.usecase.Upload(ctx.Context(), uploadDTO); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "File uploaded successfully",
	})
}
