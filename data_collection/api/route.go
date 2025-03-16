package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mntwo/tasklab/data_collection/api/data_report_api"
	"github.com/mntwo/tasklab/data_collection/api/health_check_api"
)

func Handler(route *gin.Engine) {
	route.GET("/health_check", health_check_api.HealthCheck)
	v1 := route.Group("/v1")
	{
		v1.POST("/report", data_report_api.Collect)
	}
}
