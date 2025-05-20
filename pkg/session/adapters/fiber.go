package sessionadapter

import (
	"context"
	"log"

	"github.com/Fi44er/sdmedik/backend/pkg/session"
	"github.com/gofiber/fiber/v2"
)

type fiberContext struct {
	*fiber.Ctx
}

func (f *fiberContext) Context() context.Context {
	return f.Ctx.Context()
}

func (f *fiberContext) SetContext(ctx context.Context) {
	f.Ctx.SetUserContext(ctx)
}

func (f *fiberContext) GetCookie(name string) string {
	return f.Ctx.Cookies(name)
}

func (f *fiberContext) SetCookie(name, value string, maxAge int) {
	f.Ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		HTTPOnly: true,
		MaxAge:   maxAge,
	})
}

func (f *fiberContext) SetHeader(key, value string) {
	f.Ctx.Set(key, value)
}

func FiberMiddleware(manager *session.SessionManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := &fiberContext{Ctx: c}

		ctx2, sess := manager.Start(ctx)
		c.SetUserContext(ctx2)

		ctx.SetHeader("Vary", "Cookie")
		ctx.SetHeader("Cache-Control", `no-cache="Set-Cookie"`)

		err := c.Next()

		if err := manager.Save(ctx, sess); err != nil {
			log.Printf("Failed to save session: %v", err)
		}

		return err
	}
}
