package main

import (
	"github.com/josephus-git/TCAS-simulation/internal/config"
	"github.com/josephus-git/TCAS-simulation/internal/util"
)

func resetAll(cfg *config.Config) {
	cfg.IsRunning = false
}

func restartApplication() {
	util.ResetLog()
	start()
}
