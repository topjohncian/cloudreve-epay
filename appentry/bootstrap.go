package appentry

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/topjohncian/cloudreve-pro-epay/internal/appconf"
	"go.uber.org/fx"
)

func Bootstrap() {
	opts := []fx.Option{}
	opts = append(opts, AppEntry()...)
	opts = append(opts, fx.Invoke(run))

	app := fx.New(opts...)

	app.Run() // blocks
}

func run(app *gin.Engine, conf *appconf.Config, lc fx.Lifecycle) {
	server := &http.Server{Handler: app}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			serviceLn, err := net.Listen("tcp", conf.Listen)
			if err != nil {
				return err
			}

			go func() {
				logrus.Infof("HTTP 服务器已启动，监听地址：%s", conf.Listen)
				if err := server.Serve(serviceLn); err != nil {
					if errors.Is(err, http.ErrServerClosed) {
						logrus.Infoln("HTTP 服务器已停止")
						return
					}

					logrus.WithError(err).Errorln("HTTP 服务器未预期地停止")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logrus.Infoln("HTTP 服务器停止中")
			if err := server.Shutdown(ctx); err != nil {
				logrus.WithError(err).Errorln("无法停止 HTTP 服务器")
				return err
			}

			return nil
		},
	})
}
