package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type Common struct {
	Logger    *zap.SugaredLogger
	AtomicLog *zap.AtomicLevel
}

type message struct {
	Level string `json:"level"`
}

func (s *Common) Hello(w http.ResponseWriter, r *http.Request) {
	s.Logger.Debug("debug message")
	s.Logger.Info("info message")
	s.Logger.Warn("warning message")
	s.Logger.Error("error message")
	s.Logger.DPanic("dpanic message")
}

func (s *Common) ChangeLogLevel(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var msg message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch msg.Level {
	case "info":
		s.AtomicLog.SetLevel(zap.InfoLevel)
	case "warning":
		s.AtomicLog.SetLevel(zap.WarnLevel)
	case "error":
		s.AtomicLog.SetLevel(zap.ErrorLevel)
	case "dpanic":
		s.AtomicLog.SetLevel(zap.DPanicLevel)
	case "panic":
		s.AtomicLog.SetLevel(zap.PanicLevel)
	case "fatal":
		s.AtomicLog.SetLevel(zap.FatalLevel)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
