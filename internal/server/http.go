package server

import (
	"github.com/gin-gonic/gin"
	"github.com/topjohncian/cloudreve-pro-epay/internal/appconf"
)

func CreateHttp(conf *appconf.Config) *gin.Engine {
	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	}

	r.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": conf.Listen})
	})

	return r
}
