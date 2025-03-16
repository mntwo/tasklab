package data_report_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mntwo/tasklab/event_manager"
)

func Collect(c *gin.Context) {
	m, ok := event_manager.GetEventManager("sample_task")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "event manager not found"})
		return
	}
	m.Notify(map[string]interface{}{"msg": "hello world"})
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "it's ok"})
}
