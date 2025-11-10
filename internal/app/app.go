package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/middleware"
	"github.com/Fi44er/sdmed/pkg/postgres"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
	"github.com/Fi44er/sdmed/pkg/process_manager"
	redisConnect "github.com/Fi44er/sdmed/pkg/redis"
	"github.com/Fi44er/sdmed/pkg/session"
	sessionadapter "github.com/Fi44er/sdmed/pkg/session/adapters"
	sessionstore "github.com/Fi44er/sdmed/pkg/session/store"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	app *fiber.App

	config     *config.Config
	logger     *logger.Logger
	validator  *validator.Validate
	httpConfig config.HTTPConfig

	db          *gorm.DB
	redisClient *redis.Client

	redisManager   redisConnect.IRedisManager
	sessionManager *session.SessionManager
	processManager process_manager.IProcessManager
	uow            uow.Uow

	moduleProvider *moduleProvider

	migrate   bool
	redisMode int
}

func NewApp() *App {
	migrate := flag.Bool("migrate", false, "Run database migration on startup (true/false)")
	redisMode := flag.Int("redis", 0, "Redis cache mode: 0 - no flush, 1 - selective flush, 2 - full flush")
	flag.Parse()

	return &App{
		app:       fiber.New(),
		migrate:   *migrate,
		redisMode: *redisMode,
	}
}

func (app *App) Run() error {
	app.app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:8080, http://localhost:5173, http://localhost:8080",
		AllowCredentials: true,
	}))

	app.app.Use(logger.LoggerMiddleware())
	app.app.Use(middleware.ErrHandler)

	err := app.initDeps()
	if err != nil {
		return err
	}

	if err := app.registerBackgroundProcesses(); err != nil {
		app.logger.Errorf("Failed to register background processes: %v", err)
	}

	app.processManager.StartAll()

	return app.runHttpServerWithShutdown()
}

func (app *App) initDeps() error {
	inits := []func() error{
		app.initConfig,
		app.initLogger,
		app.initDb,
		app.initRedis,
		app.initSessionManager,
		app.initValidator,
		app.initProcessManager,
		app.initModuleProvider,
		app.initRouter,
	}
	for _, init := range inits {
		err := init()
		if err != nil {
			return fmt.Errorf("✖ Failed to initialize dependencies: %s", err.Error())
		}
	}
	return nil
}

func (app *App) initConfig() error {
	if app.config == nil {
		config, err := config.LoadConfig(".")
		if err != nil {
			return fmt.Errorf("✖ Failed to load config: %s", err.Error())
		}
		app.config = config
	}

	err := config.Load(".env")
	if err != nil {
		return fmt.Errorf("✖ Failed to load config: %s", err.Error())
	}

	return nil
}

func (app *App) initDb() error {
	if app.db == nil {
		db, err := postgres.ConnectDb(app.config.PostgresUrl, app.logger)
		if err != nil {
			return err
		}
		app.db = db
		app.uow = uow.New(app.db)

		if err := postgres.Migrate(db, app.migrate, app.logger); err != nil {
			return fmt.Errorf("✖ Failed to migrate database: %s", err.Error())
		}
	}

	return nil
}

func (app *App) initRedis() error {
	if app.redisManager == nil {
		client, err := redisConnect.Connect(app.config.RedisUrl, app.logger)
		if err != nil {
			app.logger.Errorf("Failed to connect to Redis: %v", err)
			return nil
		}

		app.redisManager = redisConnect.NewRedisManger(client)
		app.redisClient = client

		if err := redisConnect.FlushRedisCache(client, app.redisMode, app.logger); err != nil {
			err = fmt.Errorf("✖ Failed to flush redis cache: %v", err)
			app.logger.Errorf("%s", err.Error())
			return err
		}
	}
	return nil
}

func (app *App) initLogger() error {
	if app.logger == nil {
		app.logger = logger.NewLogger()
	}
	return nil
}

func (app *App) initValidator() error {
	if app.validator == nil {
		app.validator = validator.New()
	}
	return nil
}

func (app *App) initSessionManager() error {
	if app.sessionManager == nil {
		app.sessionManager = session.NewSessionManager(
			sessionstore.NewRedisSessionStore(app.redisClient),
			30*time.Minute,
			1*time.Hour,
			12*time.Hour,
			"session",
		)

		app.app.Use(sessionadapter.FiberMiddleware(app.sessionManager))
	}

	return nil
}

func (app *App) initModuleProvider() error {
	err := error(nil)
	app.moduleProvider, err = NewModuleProvider(app)
	if err != nil {
		app.logger.Errorf("%s", err.Error())
		return err
	}
	return nil
}

func (app *App) runHttpServer() error {
	if app.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			app.logger.Errorf("✖ Failed to load config: %s", err.Error())
			return fmt.Errorf("✖ Failed to load config: %v", err)
		}
		app.httpConfig = cfg
	}

	app.logger.Infof("🌐 Server is running on %s", app.httpConfig.Address())
	app.logger.Info("✅ Server started successfully")
	if err := app.app.Listen(app.httpConfig.Address()); err != nil {
		app.logger.Errorf("✖ Failed to start server: %s", err.Error())
		return fmt.Errorf("✖ Failed to start server: %v", err)
	}

	return nil
}

func (app *App) runHttpServerWithShutdown() error {
	if app.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			app.logger.Errorf("✖ Failed to load config: %s", err.Error())
			return fmt.Errorf("✖ Failed to load config: %v", err)
		}
		app.httpConfig = cfg
	}

	serverErr := make(chan error, 1)

	go func() {
		app.logger.Infof("🌐 Server is running on %s", app.httpConfig.Address())
		app.logger.Info("✅ Server started successfully")
		if err := app.app.Listen(app.httpConfig.Address()); err != nil {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		app.logger.Errorf("✖ Server error: %s", err.Error())
		app.stopBackgroundProcesses()
		return err
	case sig := <-quit:
		app.logger.Infof("Received signal: %v. Shutting down...", sig)
		app.stopBackgroundProcesses()
		app.logger.Info("✅ Application stopped gracefully")
		return nil
	}
}

func (app *App) stopBackgroundProcesses() {
	if app.processManager != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := app.processManager.StopAll(ctx); err != nil {
			app.logger.Errorf("Error stopping background processes: %v", err)
		}
	}
}

func (app *App) initProcessManager() error {
	if app.processManager == nil {
		app.processManager = process_manager.NewProcessManager(app.logger)
	}
	return nil
}

func (app *App) registerBackgroundProcesses() error {
	if app.processManager == nil {
		return fmt.Errorf("process manager is not initialized")
	}

	if app.moduleProvider != nil && app.moduleProvider.fileModule != nil {
		fileCleaner := app.moduleProvider.fileModule.GetFileCleaner()
		if fileCleaner != nil {
			app.processManager.Register(fileCleaner)
			app.logger.Info("✅ FileCleaner registered in process manager")
		}
	}

	return nil
}

func (app *App) initRouter() error {
	app.app.Get("/swagger/*", swagger.HandlerDefault)
	api := app.app.Group("/api")

	app.moduleProvider.userModule.InitDelivery(api)
	app.moduleProvider.authModule.InitDelivery(api)
	app.moduleProvider.fileModule.InitDelivery(api)

	return nil
}
