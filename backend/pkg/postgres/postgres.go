package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/kotopesp/sos-kotopes/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	DB *gorm.DB
}

func New(ctx context.Context, dbURL string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	var db *gorm.DB
	var err error

	for pg.connAttempts > 0 {
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: logger.Log(),
		})
		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				return nil, fmt.Errorf("postgres - New - db.DB: %w", err)
			}
			sqlDB.SetMaxOpenConns(pg.maxPoolSize)
			sqlDB.SetConnMaxLifetime(time.Hour)

			pg.DB = db
			break
		}

		logger.Log().Debug(ctx,
			"postgres is trying to connect, attempts left: %d", pg.connAttempts,
		)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		logger.Log().Fatal(ctx, "failed to connect to database: %s", err.Error())
		return nil, err
	}

	return pg, nil
}

func (p *Postgres) Close(ctx context.Context) {
	sqlDB, err := p.DB.DB()
	if err != nil {
		logger.Log().Info(ctx, "Error getting underlying database connection: %s", err.Error())
		return
	}

	if err := sqlDB.Close(); err != nil {
		logger.Log().Info(ctx, "Error closing database connection: %s", err.Error())
	}
}
