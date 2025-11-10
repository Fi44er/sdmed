package file_http

import (
	"context"
	"io"
	"time"

	file_entity "github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	_ "github.com/Fi44er/sdmed/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IFileUsecase interface {
	UploadTemporary(ctx context.Context, file *file_entity.File, ttl time.Duration) error
}

type FileHandler struct {
	usecase IFileUsecase

	validator *validator.Validate
	logger    *logger.Logger
	converter *Converter
}

func NewFileHandler(
	usecase IFileUsecase,

	validator *validator.Validate,
	logger *logger.Logger,
) *FileHandler {
	return &FileHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
		converter: &Converter{},
	}
}

// @Summary Upload Temporary File
// @Description Uploads a temporary file to the storage with a specified time-to-live (24 hours)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} response.Response "File successfully uploaded"
// @Failure 500 {object} response.Response "Error reading file or uploading to storage"
// @Router /files/upload-temporary [post]
func (h *FileHandler) UploadTemporary(ctx *fiber.Ctx) error {
	h.logger.Debug("Pidor")
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

	entity := file_entity.File{
		Data: fileData,
	}

	if err := h.usecase.UploadTemporary(ctx.Context(), &entity, 1*time.Minute); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "File uploaded successfully",
	})
}
