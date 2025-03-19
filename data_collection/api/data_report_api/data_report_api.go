package data_report_api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mntwo/tasklab/dispatcher"
	"github.com/mntwo/tasklab/encoding/json"
	"github.com/mntwo/tasklab/internal/log"
	"go.uber.org/zap"
)

func Collect(c *gin.Context) {
	var (
		ctx = c.Request.Context()
	)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(ctx, "read body failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "read body failed"})
		return
	}
	p := json.New()
	err = p.Unmarshal(body)
	if err != nil {
		log.Error(ctx, "unmarshal body failed", zap.Error(err), zap.ByteString("body", body))
		c.JSON(http.StatusBadRequest, gin.H{"code": 2, "msg": "unmarshal body failed"})
		return
	}
	err = dispatcher.Dispatch(ctx, p)
	if err != nil {
		log.Error(ctx, "dispatch failed", zap.Error(err), zap.Any("payload", p))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 3, "msg": "dispatch failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
}
