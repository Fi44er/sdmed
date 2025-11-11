package file_http

import (
	"context"
	"io"
	"time"

	file_entity "github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/internal/module/file/pkg/utils"
	"github.com/Fi44er/sdmed/pkg/logger"
	_ "github.com/Fi44er/sdmed/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IFileUsecase interface {
	UploadTemporary(ctx context.Context, file *file_entity.File, ttl time.Duration) (string, error)
	Get(ctx context.Context, name string) (*file_entity.File, error)
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

	link, err := h.usecase.UploadTemporary(ctx.Context(), &entity, 10*time.Minute)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"url":     link,
	})
}

// Get download file by name
// @Summary Download file by name
// @Description Download a file by its name. Returns file content as attachment.
// @Tags files
// @Accept json
// @Produce application/octet-stream
// @Produce image/jpeg
// @Produce image/png
// @Produce application/pdf
// @Produce text/plain
// @Param name path string true "File name"
// @Success 200 {file} byte "File content"
// @Header 200 {string} Content-Disposition "attachment; filename=example.jpg"
// @Header 200 {string} Content-Type "MIME type of the file"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /files/{name} [get]
func (h *FileHandler) Get(ctx *fiber.Ctx) error {
	name := ctx.Params("name")

	file, err := h.usecase.Get(ctx.Context(), name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ctx.Set("Content-Disposition", "inline; filename="+file.Name)
	ctx.Set("Content-Type", "application/octet-stream")

	contentType := utils.GetContentType(file.Name)
	ctx.Set("Content-Type", contentType)

	return ctx.Send(file.Data)
}
