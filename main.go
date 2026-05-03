package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*Context) error
}

// Contains URLs used for pagination.
type Context struct {
	Previous *string
	Next     *string
}

var COMMANDS map[string]CliCommand

func main() {
	COMMANDS = map[string]CliCommand{
		"help": {"help", "Displays a help message", commandHelp},
		"map":  {"map", "Display the next 20 Pokemon map locations", commandMapNext},
		"mapb": {"mapb", "Display the previous 20 Pokemon map locations", commandMapPrevious},
		"exit": {"exit", "Exit the Pokedex", commandExit},
	}

	var mapContext *Context

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
					command.callback(mapContext)
				} else {
					fmt.Println("Unknown command")
				}
			}
		}
	}
}

func commandMapNext(context *Context) error {
	if context == nil {
		context = &Context{}
	} else if context.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}

	result, err := getLocationAreas(context.Next)
	if err != nil {
		log.Fatal(err)
	}

	context.Previous = result.Previous
	context.Next = result.Next

	if len(result.Results) == 0 {
		fmt.Println("No locations found")
		return nil
	}

	for _, location := range result.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapPrevious(context *Context) error {
	if context == nil {
		context = &Context{}
	} else if context.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	result, err := getLocationAreas(context.Previous)
	if err != nil {
		log.Fatal(err)
	}

	context.Previous = result.Previous
	context.Next = result.Next

	if len(result.Results) == 0 {
		fmt.Println("No locations found")
		return nil
	}

	for _, location := range result.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandHelp(context *Context) error {
	_, err := fmt.Println("Welcome to the Pokedex!\nUsage:")
	if err != nil {
		return err
	}
	for _, val := range COMMANDS {
		fmt.Printf("  %s: %s\n", val.name, val.description)
	}
	return nil
}

func commandExit(context *Context) error {
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
