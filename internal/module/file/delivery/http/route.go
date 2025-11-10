package file_http

import "github.com/gofiber/fiber/v2"

func (h *FileHandler) RegisterRoutes(router fiber.Router) {
	files := router.Group("/files")

	files.Post("/upload-temporary", h.UploadTemporary)
	files.Get("/:name", h.Get)
}
