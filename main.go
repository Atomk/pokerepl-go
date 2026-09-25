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
		"inspect": {"inspect", "Get information about a pokemon", commandInspect},
		"catch":   {"catch", "Try to catch a pokemon by name", commandCatch},
		"pokedex": {"pokedex", "List all caught pokemon", commandPokedex},
		"exit":    {"exit", "Exit the Pokedex", commandExit},
	}

	mapContext := NewContext(5)
	// Commands history. The most recent entry is always history[-1]
	history := NewHistory()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		// Issue: this reads when you press Enter but to detect arrows we need
		// a different approach.
		if scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			words := cleanInput(input)
			if len(words) == 0 {
				continue
			}
			if len(words) == 1 && words[0] == "\x1b[A" {
				if history.Count() > 0 {
					fmt.Printf("\r%s", history.Previous())
				}
			}
			if len(words) > 0 {
				commandName := words[0]
				arguments := words[1:]
				command, ok := COMMANDS[commandName]
				if ok {
					history.Add(input)

					err := command.callback(mapContext, arguments)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println()
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
