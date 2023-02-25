package server

import (
	"html/template"
	"io/fs"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/topjohncian/cloudreve-pro-epay/internal/appconf"
)

func CreateHttp(conf *appconf.Config, templateFS fs.FS) *gin.Engine {
	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	}

	tmpl := lo.Must(template.ParseFS(templateFS, "templates/*.tmpl"))
	r.SetHTMLTemplate(tmpl)

	r.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": conf.Listen})
	})

	return r
}
