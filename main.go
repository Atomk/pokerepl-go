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
	callback    func(*Context, []string) error
}

var COMMANDS map[string]CliCommand

func main() {
	COMMANDS = map[string]CliCommand{
		"help":    {"help", "Displays a help message", commandHelp},
		"map":     {"map", "Display the next 20 Pokemon map locations", commandMapNext},
		"mapb":    {"mapb", "Display the previous 20 Pokemon map locations", commandMapPrevious},
		"explore": {"explore", "Explore the area passed as an argument", commandExplore},
		"catch":   {"catch", "Try to catch a pokemon by name", commandCatch},
		"exit":    {"exit", "Exit the Pokedex", commandExit},
	}

	mapContext := NewContext(5)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			words := cleanInput(input)
			if len(words) > 0 {
				commandName := words[0]
				arguments := words[1:]
				command, ok := COMMANDS[commandName]
				if ok {
					command.callback(mapContext, arguments)
				} else {
					fmt.Println("Unknown command")
				}
			}
		}
	}
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)
	return words
}
