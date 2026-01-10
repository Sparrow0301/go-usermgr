package svc

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"usermgmt/internal/config"
	"usermgmt/internal/middleware"
	"usermgmt/internal/model"
)

// ServiceContext wires together shared resources that handlers and logic layers rely on.
type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	Validator      *validator.Validate
	AuthMiddleware rest.Middleware
	RoleGuard      func(roles ...string) rest.Middleware
}

// NewServiceContext builds the service context with DB, validator and middlewares.
func NewServiceContext(c config.Config) *ServiceContext {
	db := mustInitDB(c)
	validate := validator.New(validator.WithRequiredStructEnabled())

	ctx := &ServiceContext{
		Config:    c,
		DB:        db,
		Validator: validate,
	}
	ctx.AuthMiddleware = middleware.NewAuthMiddleware(c.JWT.AccessSecret).Handle
	ctx.RoleGuard = func(roles ...string) rest.Middleware {
		return middleware.NewRoleGuard(roles...)
	}
	return ctx
}

// AutoMigrate ensures schema is created or updated.
func (s *ServiceContext) AutoMigrate() error {
	return s.DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
	)
}

// mustInitDB establishes the GORM connection and tunes the connection pool.
func mustInitDB(c config.Config) *gorm.DB {
	gormLogger := logger.New(
		log.New(log.Writer(), "GORM", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(c.Database.DSN), &gorm.Config{Logger: gormLogger})
	if err != nil {
		logx.Errorf("failed to connect database: %v", err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	maxIdle := c.Database.MaxIdleConns
	if maxIdle <= 0 {
		// Default to sane pool numbers even when config omits them.
		maxIdle = 10
	}
	maxOpen := c.Database.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 30
	}
	connLifetime := c.Database.ConnMaxLifetime
	if connLifetime <= 0 {
		connLifetime = time.Hour
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(connLifetime)

	return db
}
