package app

import (
	"github.com/mntwo/tasklab/internal/application"
	"github.com/mntwo/tasklab/internal/db"
)

func NewDatabaseApp() application.Application {
	return db.New("database")
}
