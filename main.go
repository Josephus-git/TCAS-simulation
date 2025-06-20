package main

import (
	"os"
)

func main() {
	// reset all log files
	resetLog()
	startR()
}

func resetLog() {
	filesToDelete := []string{
		"airportDetails.txt",
		"airPlaneDetails.txt",
		"flightDetails.txt",
	}

	for _, filePath := range filesToDelete {
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			continue
		}
		os.Remove(filePath)
	}
}
