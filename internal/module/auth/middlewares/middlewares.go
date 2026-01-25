package auth_middlewares

import (
	"fmt"

	"github.com/Fi44er/sdmed/internal/config"
	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/go-viper/mapstructure/v2"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	logger *logger.Logger
	config *config.Config
}

func NewAuthMiddleware(logger *logger.Logger, config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,
		config: config,
	}
}

type ITokenService interface {
	ValidateToken(token, publicKey string) (*auth_entity.TokenDetails, error)
}

func (m *AuthMiddleware) Guest() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess := session.FromFiberContext(ctx)
		if sess == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session not found",
			})
		}

		sessionData, ok := sess.Get("session_info").(map[string]any)
		if !ok || sessionData == nil {
			data := auth_entity.ActiveSession{
				Roles: []string{"guest"},
			}

			// проверять авторизован ли пользователь
			// если не авторизирован то проверяем user_id
			// если есть то получаем из бд
			// если нет то создаем теневой аккаунт и записываем его user_id в долгоживущую куки и заполняем новую сессию данными

			var newData map[string]any
			if err := mapstructure.Decode(data, &newData); err != nil {
				return fmt.Errorf("failed to encode session data for guest user: %w", err)
			}
			sess.Put("session_info", newData)
		}

		return ctx.Next()
	}
}
