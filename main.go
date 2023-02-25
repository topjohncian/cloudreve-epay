package main

import (
	"embed"
	"flag"
	"io/fs"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/topjohncian/cloudreve-pro-epay/appentry"
)

//go:embed templates/*
var templateFS embed.FS

var isEject bool

func init() {
	flag.BoolVar(&isEject, "eject", false, "导出模板文件")
	flag.Parse()
}

func main() {
	if isEject {
		appentry.Eject(templateFS)
		return
	}

	var tmplFS fs.FS
	if appentry.Exists("custom") {
		logrus.Infoln("使用自定义模板文件")
		tmplFS = os.DirFS("custom")
	} else {
		tmplFS = templateFS
	}

	appentry.Bootstrap(tmplFS)
}
