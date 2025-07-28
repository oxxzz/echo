package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
	"v5/engine"
	"v5/internal/db"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// Initialize the configuration
	viper.AutomaticEnv()
	viper.SetEnvPrefix("OX_")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(".")
	viper.SetConfigName("cfg")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("error reading config file: %v", err)
	}

	if viper.GetBool("app.debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	fWriter := lumberjack.Logger{
		Filename: viper.GetString("log.path"),
		MaxSize:  200,
		MaxAge:   7,
		Compress: true,
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(&fWriter, os.Stdout))
	logrus.Info("config loaded successfully")

	// Initialize the database
	if viper.GetBool("db.mysql.enabled") {
		if err := db.SetupMySQL(); err != nil {
			logrus.Fatalf("error initializing MySQL: %v", err)
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)

	ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer stop()

	eng := engine.New()
	go func() {
		listen := viper.GetString("app.listen")
		if listen == "" {
			listen = ":8000"
		}

		if err := eng.Start(listen); err != nil && err != http.ErrServerClosed {
			eng.Logger.Fatal("shutting down the server")
		}
	}()

	logrus.Tracef("server starting")
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	if err := eng.Shutdown(ctx); err != nil {
		eng.Logger.Fatal(err)
	}

	logrus.Tracef("server stopped")
}
