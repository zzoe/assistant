package app

import (
	"github.com/zzoe/assistant/app/attendance"
	"github.com/zzoe/assistant/cfg"
	"go.uber.org/zap"
)

var (
	log = cfg.Log
)

func Launch() {
	log.Info("Begin launch app")

	code, err := attendance.Run()
	if err != nil {
		log.Error("window err", zap.Int("code", code), zap.Error(err))
	}
}
