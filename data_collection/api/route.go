package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mntwo/tasklab/data_collection/api/health_check_api"
)

func Handler(route *gin.Engine) {
	v1 := route.Group("/v1")
	{
		v1.GET("/health_check", health_check_api.HealthCheck)
	}
}
