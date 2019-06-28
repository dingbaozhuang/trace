package service

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yumimobi/trace/models"
)

func defaultHandler(c *gin.Context) {
	c.String(http.StatusOK, "trace service\n"+"time:201905091020"+
		`{"bid_id":"0bts0B1HqucD3wSCT42khQho1EgJ3p","mac":"10:44:00:e4:9d:e1","language":"zh_CN","telMake":"HONOR","brand":"","telModel":"BKL-AL20","carrier":0,"terminalType":0,"netEnv":3,"screenWidth":720,"screenHeigh":1440,"screenDensity":"360","lat":"0","lng":"0","os":0,"osVersion":"8.0","sdkVersion":"android3.4.0","adType":2,"screenMode":1,"adWidth":640,"adHeigh":100,"clientIp":"119.0.201.166","ua":"Mozilla/5.0 (Linux; Android 8.0.0; BKL-AL20 Build/HUAWEIBKL-AL20; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/68.0.3440.91 Mobile Safari/537.36","appKey":"8dc5cbb16ebb618d786da5f1f4573bc6","sspId":5,"idfa":"","idfv":"","openudid":"","imsi":"460038640888157","imei":"868341031235756","androidId":"c154b2963d66cec9","android_adid":"","partner_id":10000,"tailFlag":0,"zApp_id":"8dc5cbb16ebb618d786da5f1f4573bc6","appVersion":"5.3.0","packageName":"com.brianbaek.game.popstar","appName":"国内-安卓-消灭星星官方正版","inventory_types":[1,2,5,4,7,9],"isAdxToNative":false,"sspAppIdKey":"847726511-DA5687-C740-D98B-7410B1E7A","sspAppPlaceIdKey":"847726511pjpu7f","sspAppSecretKey":"BE541AA52E52881BCA41326758CAED2E","slot_id":"d2twem9v"}`)
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
