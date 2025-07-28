package engine

import (
	"v5/engine/hooks"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
)

func New() *echo.Echo {
	eng := echo.New()
	eng.HideBanner = true
	eng.HidePort = true
	eng.JSONSerializer = &hooks.SonicJSONSerializer{}

	if viper.GetBool("log.access.enabled") {
		eng.Use(middleware.Logger())
	}

	if viper.GetString("log.access.path") != "" {
		writers := &lumberjack.Logger{
			Filename:   "storage/logs/access.logs",
			MaxSize:    200,
			MaxAge:     7,
			MaxBackups: 3,
			Compress:   true,
		}
		eng.Logger.SetOutput(writers)
	}

	if viper.GetBool("app.debug") {
		eng.Logger.SetLevel(log.DEBUG)
	}

	eng.Use(middleware.Gzip())
	eng.Use(middleware.Recover())
	eng.Use(middleware.RequestID())

	setupRoutes(eng)
	return eng
}
