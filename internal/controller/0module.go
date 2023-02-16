package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	"github.com/topjohncian/cloudreve-pro-epay/internal/appconf"
	"github.com/topjohncian/cloudreve-pro-epay/internal/cache"
	"go.uber.org/fx"
)

type CloudrevePayController struct {
	fx.In

	Conf   *appconf.Config
	Cache  cache.Driver
	Client *req.Client
}

func RegisterControllers(c CloudrevePayController, r *gin.Engine) {
	r.POST("/cloudreve/purchase", c.BearerAuthMiddleware(), c.Purchase)
	r.GET("/purchase/:id", c.PurchasePage)
	r.GET("/notify/:id", c.Notify)
	r.GET("/return/:id", c.Return)
}

func Module() fx.Option {
	return fx.Module("controller", fx.Invoke(RegisterControllers))
}
