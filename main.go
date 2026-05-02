package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CliCommand struct {
	name        string
	description string
	callback    func() error
}

var COMMANDS map[string]CliCommand

func main() {
	COMMANDS = map[string]CliCommand{
		"help": {"help", "Displays a help message", commandHelp},
		"exit": {"exit", "Exit the Pokedex", commandExit},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			words := cleanInput(input)
			if len(words) > 0 {
				commandName := words[0]
				command, ok := COMMANDS[commandName]
				if ok {
					command.callback()
				} else {
					fmt.Println("Unknown command")
				}
			}
		}
	}
}

func commandHelp() error {
	_, err := fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	if err != nil {
		return err
	}
	for _, val := range COMMANDS {
		fmt.Printf("%s: %s\n", val.name, val.description)
	}
	return nil
}

func commandExit() error {
	_, err := fmt.Println("Closing the Pokedex... Goodbye!")
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)
	return words
}
