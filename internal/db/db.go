package db

import (
	"context"
	"fmt"

	"github.com/mntwo/tasklab/internal/application"
	"github.com/mntwo/tasklab/internal/config"
	"github.com/mntwo/tasklab/internal/log"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type db struct {
	Name   string
	stopCh chan struct{}
}

func New(name string) application.Application {
	return &db{
		Name:   name,
		stopCh: make(chan struct{}),
	}
}

func (d *db) Start() error {
	dbs = make(map[string]*gorm.DB)

	for _, dbConfig := range config.GetPostgres() {
		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Database,
		)

		db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err != nil {
			log.Fatal(context.Background(), "Error opening database", zap.String("name", dbConfig.Name), zap.Error(err))
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(context.Background(), "Error getting database instance", zap.String("name", dbConfig.Name), zap.Error(err))
		}

		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
		sqlDB.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime)

		if err = sqlDB.Ping(); err != nil {
			log.Fatal(context.Background(), "Error connecting to the database", zap.String("name", dbConfig.Name), zap.Error(err))
		}
		dbs[dbConfig.Name] = db
	}
	<-d.stopCh
	return application.ErrApplicationClosed
}

func (d *db) Stop() error {
	for _, db := range dbs {
		sqlDB, err := db.DB()
		if err == nil && sqlDB != nil {
			sqlDB.Close()
		}
	}
	close(d.stopCh)
	return nil
}

func (d *db) GetName() string {
	return d.Name
}
