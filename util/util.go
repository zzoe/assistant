package util

import (
	"github.com/zzoe/assistant/cfg"
	"go.uber.org/zap"
)

var (
	log = cfg.Log
)

func Warn(err error) {
	if err != nil {
		log.Warn("WARN", zap.Error(err))
	}
}
