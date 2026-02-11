package middlewares

import (
	"context"
	"fmt"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	accessmanager_service "github.com/Fi44er/sdmed/internal/module/auth/usecase/access_manager"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/go-viper/mapstructure/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IShadowUserService interface {
	CreateShadowUser(ctx context.Context) (*auth_entity.User, error)
}

type ISessionRepository interface {
	GetSessionInfo(ctx context.Context) (*auth_entity.ActiveSession, error)
	PutSessionInfo(ctx context.Context, sessionInfo *auth_entity.ActiveSession) error
}

type IUserSessionRepository interface {
	Create(ctx context.Context, session *auth_entity.UserSession) error
	Get(ctx context.Context, id string) (*auth_entity.UserSession, error)
	UpdateLastUsed(ctx context.Context, id string) error
}

func InjectManager(am *accessmanager_service.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(ManagerKey, am)
		return c.Next()
	}
}

const ManagerKey = "accessManager"

func ShadowSessionMiddleware(
	shadowUserService IShadowUserService,
	sessionRepository ISessionRepository,
	userSessionRepository IUserSessionRepository,
	logger *logger.Logger,
	config *config.Config,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем сессию из контекста
		sess := session.FromFiberContext(c)
		if sess == nil {
			logger.Debug("Session not found in context, skipping shadow session creation")
			return c.Next()
		}

		ctx := context.WithValue(c.Context(), "session", *sess)

		// КРИТИЧЕСКИ ВАЖНО: Проверяем есть ли уже информация в сессии
		sessionInfo, err := sessionRepository.GetSessionInfo(ctx)

		logger.Warnf("Sess ebanay %+v", sessionInfo)

		// Если сессия существует - только обновляем время последней активности
		if err == nil && sessionInfo != nil {
			logger.Debugf("Session already exists for device: %s", sessionInfo.DeviceID)

			// Обновляем время последней активности в БД (асинхронно, не блокируем запрос)
			if sessionInfo.DeviceID != "" {
				go func() {
					updateCtx := context.Background()
					if updateErr := userSessionRepository.UpdateLastUsed(updateCtx, sessionInfo.DeviceID); updateErr != nil {
						logger.Errorf("Failed to update last used time: %v", updateErr)
					}
				}()
			}

			return c.Next()
		}

		// Если сессия НЕ существует - создаем shadow session
		logger.Debug("No session found, creating shadow session")

		// Создаем shadow user
		shadowUser, err := shadowUserService.CreateShadowUser(ctx)
		if err != nil {
			logger.Errorf("Failed to create shadow user: %v", err)
			// НЕ возвращаем ошибку, просто продолжаем без shadow session
			return c.Next()
		}

		// Генерируем уникальный ID устройства
		deviceID := uuid.New().String()

		// Получаем информацию о клиенте
		ipAddress := c.IP()
		userAgent := c.Get("User-Agent")

		// Создаем активную сессию
		now := time.Now()
		expiresAt := now.Add(config.RefreshTokenExpiresIn)

		activeSession := &auth_entity.ActiveSession{
			UserID:    shadowUser.ID,
			DeviceID:  deviceID,
			UserRoles: []string{"guest"},
			IsShadow:  true,
			CreatedAt: now,
			ExpiresAt: expiresAt,
			IPAddress: ipAddress,
			UserAgent: userAgent,
		}

		// Сохраняем в Redis - это ОБЯЗАТЕЛЬНАЯ операция
		if err := sessionRepository.PutSessionInfo(ctx, activeSession); err != nil {
			logger.Errorf("Failed to save session info to Redis: %v", err)
			// Если не удалось сохранить в Redis, не создаем запись в БД
			return c.Next()
		}

		// Сохраняем в БД для управления устройствами
		dbSession := &auth_entity.UserSession{
			ID:         deviceID,
			UserID:     shadowUser.ID,
			DeviceName: parseDeviceName(userAgent),
			LastIP:     ipAddress,
			UserAgent:  userAgent,
			LastUsedAt: now,
			ExpiresAt:  expiresAt,
			IsRevoked:  false,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		if err := userSessionRepository.Create(ctx, dbSession); err != nil {
			logger.Errorf("Failed to create user session in DB: %v", err)
			// Даже если не удалось сохранить в БД, сессия в Redis уже есть
		} else {
			logger.Infof("Created shadow session for user %s with device %s", shadowUser.ID, deviceID)
		}

		return c.Next()
	}
}

func Guest() fiber.Handler {
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
				UserRoles: []string{"guest"},
				IsShadow:  true,
			}

			var newData map[string]any
			if err := mapstructure.Decode(data, &newData); err != nil {
				return fmt.Errorf("failed to encode session data for guest user: %w", err)
			}
			sess.Put("session_info", newData)
		}

		return ctx.Next()
	}
}

func Authorize(obj, act string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		val := ctx.Locals(ManagerKey)
		if val == nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Access manager not found",
			})
		}

		am := val.(*accessmanager_service.Manager)

		sess := session.FromFiberContext(ctx)
		if sess == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session not found",
			})
		}

		sessionData, ok := sess.Get("session_info").(map[string]any)
		if !ok {
			return ctx.Status(403).JSON(fiber.Map{
				"error": "Invalid session data",
			})
		}

		var userSession auth_entity.ActiveSession
		if err := mapstructure.Decode(sessionData, &userSession); err != nil {
			return fmt.Errorf("failed to decode session data: %v", err)
		}

		hasAccess := false
		if userSession.UserRoles != nil {
			for _, role := range userSession.UserRoles {
				ok, _ := am.Enforcer.Enforce(role, obj, act)
				if ok {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}

		return ctx.Next()
	}
}

func RequireAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess := session.FromFiberContext(ctx)
		if sess == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session not found",
			})
		}

		sessionData, ok := sess.Get("session_info").(map[string]any)
		if !ok || sessionData == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Not authenticated",
			})
		}

		var userSession auth_entity.ActiveSession
		if err := mapstructure.Decode(sessionData, &userSession); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid session data",
			})
		}

		// Проверяем, что пользователь не shadow
		if userSession.IsShadow {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Authentication required",
				"message": "Please sign in to access this resource",
			})
		}

		return ctx.Next()
	}
}

func parseDeviceName(userAgent string) string {
	// Простая реализация - можно улучшить с помощью библиотеки для парсинга User-Agent
	if len(userAgent) > 50 {
		return userAgent[:50] + "..."
	}
	if userAgent == "" {
		return "Unknown Device"
	}
	return userAgent
}
