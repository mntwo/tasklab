package main

import (
	"github.com/mntwo/tasklab/internal/app"
	_ "github.com/mntwo/tasklab/internal/config"
)

func main() {
	dataCollectApp := app.NewDataCollectionApp()
	genEventApp := app.NewGenEventApp()
	app.Run(dataCollectApp, genEventApp)
}
