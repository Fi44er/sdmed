package app

import (
	auth_module "github.com/Fi44er/sdmedik/backend/internal/module/auth"
	notification_module "github.com/Fi44er/sdmedik/backend/internal/module/notification"
	user_module "github.com/Fi44er/sdmedik/backend/internal/module/user"
)

type moduleProvider struct {
	app *App

	userModule         *user_module.UserModule
	notificationModule *notification_module.NotificationModule
	authModule         *auth_module.AuthModule
}

func NewModuleProvider(app *App) (*moduleProvider, error) {
	provider := &moduleProvider{
		app: app,
	}

	err := provider.initDeps()
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (p *moduleProvider) initDeps() error {
	inits := []func() error{
		p.UserModule,
		p.NotificationModule,
		p.AuthModule,
	}
	for _, init := range inits {
		err := init()
		if err != nil {
			p.app.logger.Errorf("%s", "✖ Failed to initialize module: "+err.Error())
			return err
		}
	}
	return nil
}

func (p *moduleProvider) UserModule() error {
	p.userModule = user_module.NewUserModule(p.app.logger, p.app.validator, p.app.db)
	p.userModule.Init()
	return nil
}

func (p *moduleProvider) NotificationModule() error {
	p.notificationModule = notification_module.NewNotificationModule(p.app.logger, p.app.config)
	p.notificationModule.Init()
	return nil
}

func (p *moduleProvider) AuthModule() error {
	p.authModule = auth_module.NewAuthModule(
		p.app.logger,
		p.app.validator,
		p.app.db,
		p.app.redisManager,
		p.app.config,
		p.userModule.UserUsecase,
		p.notificationModule.Service,
	)
	p.authModule.Init()
	return nil
}
