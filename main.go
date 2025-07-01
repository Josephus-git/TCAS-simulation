package main

import (
	"bufio"
	"fmt"
	"internal/config"
	"os"
)

func main() {
	resetLog()
	start()
}

// start initializes the TCAS simulator, loads configurations, and enters a continuous command-line interaction loop.
func start() {
	scanner := bufio.NewScanner(os.Stdin)
	initialize := &config.Config{}

	getNumberPlanes(initialize)
	initializeAirports(initialize)

	for i := 0; ; i++ {
		fmt.Print("TCAS-simulator > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		argument2 := ""
		if len(input) > 1 {
			argument2 = input[1]
		}

		cmd, ok := getCommand(api, argument2)[input[0]]
		if !ok {
			fmt.Println("Unknown command, type <help> for usage")
			continue
		}
		err := cmd.callback()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		println("")

	}
}
