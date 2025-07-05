package main

import (
	"github.com/josephus-git/TCAS-simulation/internal/config"
	"github.com/josephus-git/TCAS-simulation/internal/util"
)

// resetAll sets the `IsRunning` and `DifferentAltitudes` flags in the provided configuration to false.
func resetAll(cfg *config.Config) {
	cfg.IsRunning = false
	cfg.DifferentAltitudes = false
}

// restartApplication resets the application's log and then initiates the application startup sequence.
func restartApplication() {
	util.ResetLog()
	start()
}
