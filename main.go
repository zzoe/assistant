package main

import (
	"github.com/spf13/viper"
	"github.com/zzoe/assistant/app"
	_ "github.com/zzoe/assistant/app"
	"github.com/zzoe/assistant/cfg"
	"go.uber.org/zap"
)

var (
	log = cfg.Log
)

func main() {
	log.Debug("main begin")
	defer func() {
		if err := log.Sync(); err != nil {
			panic(err)
		}

	}()

	order := viper.GetStringMapString("excel.order")
	log.Info("map", zap.Any("order", order), zap.String("c", order["c"]))

	app.Launch()
}
