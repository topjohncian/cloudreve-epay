package appconf

import (
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../..")
)

func Parse() (*Config, error) {
	err := godotenv.Load(filepath.Join(Root, ".env"))

	if err != nil {
		logrus.WithError(err).Fatalln("无法加载 .env 文件")
		return nil, err
	}

	var config Config
	err = envconfig.Process("cr_epay", &config)

	if err != nil {
		envconfig.Usage("cr_epay", &config)
		logrus.WithError(err).Fatalln("无法加载配置")
		return nil, err
	}

	return &config, nil
}
