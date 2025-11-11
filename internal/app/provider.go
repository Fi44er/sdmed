package app

import (
	auth_module "github.com/Fi44er/sdmed/internal/module/auth"
	file_module "github.com/Fi44er/sdmed/internal/module/file"
	notification_module "github.com/Fi44er/sdmed/internal/module/notification"
	product_module "github.com/Fi44er/sdmed/internal/module/product"
	user_module "github.com/Fi44er/sdmed/internal/module/user"
)

type moduleProvider struct {
	app *App

	userModule         *user_module.UserModule
	notificationModule *notification_module.NotificationModule
	authModule         *auth_module.AuthModule
	fileModule         *file_module.FileModule
	productModule      *product_module.ProductModule
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
		p.FileModule,
		p.ProductModule,
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
		p.userModule.GetUserUsecase(),
		p.notificationModule.GetNotificationService(),
	)
	p.authModule.Init()
	return nil
}

func (p *moduleProvider) FileModule() error {
	p.fileModule = file_module.NewFileModule(
		p.app.logger,
		p.app.validator,
		p.app.db,
		p.app.config,
		p.app.uow,
	)
	p.fileModule.Init()
	return nil
}

func (p *moduleProvider) ProductModule() error {
	p.productModule = product_module.NewProductModule(
		p.app.logger,
		p.app.validator,
		p.app.db,
		p.app.uow,
		p.fileModule.GetFileService(),
	)
	p.productModule.Init()
	return nil
}
