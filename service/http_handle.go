package service

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/models"
)

const (
	VERSION = "1.0.0"
)

func defaultHandler(c *gin.Context) {
	conf := config.Conf
	c.String(http.StatusOK, "yumimobi trace service\n"+"version:"+VERSION+"\n"+conf.Server.HTTP.Address+":"+conf.Server.HTTP.Port+`/trace`)
}

func heartbeatHandler(c *gin.Context) {
	resp := &models.HTTPResponse{
		Code: models.StatusSuccess,
	}
	c.JSON(http.StatusOK, resp)
}

func shutdownHandler(c *gin.Context) {
	resp := &models.HTTPResponse{
		Code: models.StatusSuccess,
	}
	c.JSON(http.StatusOK, resp)
	os.Exit(0)
}

func updateTimerHandler(c *gin.Context) {
	ResetTimer()

	resp := &models.HTTPResponse{
		Code: models.StatusSuccess,
	}
	c.JSON(http.StatusOK, resp)
}

func traceHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
