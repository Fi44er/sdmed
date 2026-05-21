package app

import (
	auth_module "github.com/Fi44er/sdmed/internal/module/auth"
	cart_module "github.com/Fi44er/sdmed/internal/module/cart"
	file_module "github.com/Fi44er/sdmed/internal/module/file"
	notification_module "github.com/Fi44er/sdmed/internal/module/notification"
	product_module "github.com/Fi44er/sdmed/internal/module/product"
	scraper_module "github.com/Fi44er/sdmed/internal/module/scraper"
	tru_module "github.com/Fi44er/sdmed/internal/module/tru"
	user_module "github.com/Fi44er/sdmed/internal/module/user"
	order_module "github.com/Fi44er/sdmed/internal/module/order"
	chat_module "github.com/Fi44er/sdmed/internal/module/chat"
)

type moduleProvider struct {
	app *App

	userModule         *user_module.UserModule
	notificationModule *notification_module.NotificationModule
	authModule         *auth_module.AuthModule
	fileModule         *file_module.FileModule
	productModule      *product_module.ProductModule
	scraperModule      *scraper_module.ScraperModule
	truModule          *tru_module.TruModule
	cartModule         *cart_module.CartModule
	orderModule        *order_module.OrderModule
	chatModule         *chat_module.ChatModule
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
		p.FileModule,
		p.ProductModule,
		p.TruModule,
		p.ScraperModule,
		p.CartModule,
		p.OrderModule,
		p.ChatModule,
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
	p.notificationModule = notification_module.NewNotificationModule(p.app.logger, p.app.config, p.app.db)
	p.notificationModule.Init()
	return nil
}

func (p *moduleProvider) ScraperModule() error {
	p.scraperModule = scraper_module.NewScraperModule(
		p.app.logger,
		p.app.validator,
		p.app.db,
		p.productModule.GetProductUsecase(),
		p.truModule.GetTRUCodeUsecase(),
	)
	p.scraperModule.Init()
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
		p.userModule.GetRoleUsecase(),
		p.notificationModule.GetNotificationService(),
		p.cartModule.GetCartMigrator(),
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
		p.app.config,
		p.app.redisManager,
	)
	p.productModule.Init()
	return nil
}

func (p *moduleProvider) TruModule() error {
	p.truModule = tru_module.NewTruModule(
		p.app.logger,
		p.app.db,
	)
	p.truModule.Init()
	return nil
}

func (p *moduleProvider) CartModule() error {
	p.cartModule = cart_module.NewCartModule(
		p.app.db,
		p.app.logger,
		p.app.validator,
		p.app.uow,
		p.productModule.GetProductUsecase(),
		p.truModule.GetTRUCodeUsecase(),
	)
	p.cartModule.Init()
	return nil
}

func (p *moduleProvider) OrderModule() error {
	p.orderModule = order_module.NewOrderModule(
		p.app.db,
		p.app.logger,
		p.app.validator,
		p.app.config,
		p.cartModule.GetUsecase(),
		p.productModule.GetProductUsecase(),
		p.truModule.GetTRUCodeUsecase(),
		p.notificationModule.GetNotificationService(),
	)
	p.orderModule.Init()
	return nil
}

func (p *moduleProvider) ChatModule() error {
	p.chatModule = chat_module.NewChatModule(
		p.app.logger,
		p.app.db,
		p.app.validator,
		p.app.redisClient,
		p.userModule.GetUserUsecase(),
		p.notificationModule.GetNotificationService(),
	)
	p.chatModule.Init()
	return nil
}
