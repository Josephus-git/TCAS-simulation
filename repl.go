package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type apiConfig struct {
	noOfAirplanes int
	listAirports  []airport
}

func startR() {
	scanner := bufio.NewScanner(os.Stdin)
	api := &apiConfig{}

	getNumberPlanes(api)
	initializeAirports(api)

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

func getNumberPlanes(conf *apiConfig) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to TCAS-simulator")
	notValidInput := true

	for i := 0; notValidInput; i++ {

		fmt.Print("Input the number of planes for the simulation > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		num, err := strconv.Atoi(input[0])
		if err != nil {
			fmt.Println("Please input a valid integer")
			continue
		}
		if num < 2 {
			fmt.Println("Please input a valid integer greater than 1")
			continue
		}

		conf.noOfAirplanes = num
		notValidInput = false
	}

}

func cleanInput(text string) []string {
	words := []string{}
	sText := strings.Split(strings.TrimSpace(text), " ")
	for _, word := range sText {
		if len(word) != 0 {
			words = append(words, strings.ToLower(word))
		}
	}
	return words
}
