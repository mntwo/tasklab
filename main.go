package main

import (
	"github.com/mntwo/tasklab/internal/app"
	_ "github.com/mntwo/tasklab/internal/config"
)

func main() {
	db := app.NewDatabaseApp()
	dataCollectApp := app.NewDataCollectionApp()
	genEventApp := app.NewGenEventApp()
	app.Run(db, dataCollectApp, genEventApp)
}
