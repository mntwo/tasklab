package main

import (
	"github.com/mntwo/tasklab/internal/app"
	_ "github.com/mntwo/tasklab/internal/config"
)

func main() {
	dataCollectApi := app.NewDataCollectionApi()
	app.Run(dataCollectApi)
}
