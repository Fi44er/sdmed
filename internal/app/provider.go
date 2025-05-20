package app

import (
	notification_module "github.com/Fi44er/sdmedik/backend/internal/module/notification"
	user_module "github.com/Fi44er/sdmedik/backend/internal/module/user"
)

type moduleProvider struct {
	app *App

	userModule         *user_module.UserModule
	notificationModule *notification_module.NotificationModule
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
