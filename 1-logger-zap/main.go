package main

import (
	"go-playground/1-logger-zap/handlers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()

	c := &handlers.Common{
		Logger:    logger.Sugar(),
		AtomicLog: &atom,
	}

	router := mux.NewRouter()
	router.HandleFunc("/hello", c.Hello)
	router.HandleFunc("/changelevel", c.ChangeLogLevel)

	http.ListenAndServe(":8080", router)
}
