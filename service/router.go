package service

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Router(r *gin.Engine) {
	// 加载 /trace 模板
	r.LoadHTMLGlob("views/*.html")
	r.GET("/", defaultHandler)
	r.GET("/trace", traceHandler)
	r.GET("/server_api/:role/:date", serverAPIHTTPHandler)
}

func (c *Client) Router(r *gin.Engine) {
	r.GET("/heartbeat", heartbeatHandler)
	r.GET("/shutdown", shutdownHandler)
	r.GET("/updatetimer", updateTimerHandler)
}
