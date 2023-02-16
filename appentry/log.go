package appentry

import (
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/topjohncian/cloudreve-pro-epay/internal/appconf"
	"go.uber.org/fx/fxevent"
)

func init() {
	if os.Getenv("VSCODE_CWD") != "" {
		color.NoColor = false
	}
}

type fxLogger struct {
	l *logrus.Entry
}

var _ io.Writer = &fxLogger{}

func FxLogger(log *logrus.Logger) fxevent.Logger {
	logger := fxLogger{
		l: log.WithField("event", "fx.init"),
	}
	return &fxevent.ConsoleLogger{
		W: logger,
	}
}

func (l fxLogger) Write(p []byte) (n int, err error) {
	// from https://github.com/rs/zerolog/blob/a9a8199d2dd3578d37e459618515f34b5e917f8d/log.go#L435-L441
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		// Trim CR added by stdlog.
		p = p[0 : n-1]
	}
	l.l.Info(string(p))
	return
}

func Log(conf *appconf.Config) *logrus.Logger {
	logger := logrus.StandardLogger()
	logger.SetOutput(os.Stdout)

	if conf.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
