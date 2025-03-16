package app

import (
	"github.com/LyricTian/gzip"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mntwo/tasklab/data_collection/api"
	"github.com/mntwo/tasklab/internal/application"
	"github.com/mntwo/tasklab/internal/application/http_application"
)

func NewDataCollectionApp() application.Application {
	router := gin.New()

	router.SecureJsonPrefix("")
	router.MaxMultipartMemory = (1 << 20) * 30

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Recovery())

	api.Handler(router)
	return http_application.New("data_collection_api", http_application.WithHandler(router))
}
