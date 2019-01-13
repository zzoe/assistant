package main

import (
	"github.com/zzoe/assistant/app"
	_ "github.com/zzoe/assistant/app"
	"github.com/zzoe/assistant/cfg"
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

	app.Launch()
}
