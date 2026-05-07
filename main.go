package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pokecache "github.com/atomk/pokedexcli/internal"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*Context, []string) error
}

// Contains URLs used for pagination.
// If both previous and Next are null, no request was made to get location areas.
type Context struct {
	cache    *pokecache.Cache
	Previous *string
	Next     *string
}

var COMMANDS map[string]CliCommand

func main() {
	COMMANDS = map[string]CliCommand{
		"help":    {"help", "Displays a help message", commandHelp},
		"map":     {"map", "Display the next 20 Pokemon map locations", commandMapNext},
		"mapb":    {"mapb", "Display the previous 20 Pokemon map locations", commandMapPrevious},
		"explore": {"explore", "Explore the area passed as an argument", commandExplore},
		"exit":    {"exit", "Exit the Pokedex", commandExit},
	}

	mapContext := &Context{
		cache: pokecache.NewCache(5 * time.Minute),
	}

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

func commandMapNext(context *Context, args []string) error {
	if context.Previous != nil && context.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}

	result, err := getLocationAreas(context.Next, context.cache)
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

func commandMapPrevious(context *Context, args []string) error {
	if context.Next != nil && context.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	result, err := getLocationAreas(context.Previous, context.cache)
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

func commandExplore(context *Context, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("command `explore` requires exactly one argument")
	}

	idOrName := strings.TrimSpace(args[0])
	if len(idOrName) == 0 {
		return fmt.Errorf("provided argument is an empty string")
	}

	result, err := getLocationArea(idOrName, context.cache)
	if err != nil {
		log.Fatal(err)
	}

	if len(result.PokemonEncounters) == 0 {
		fmt.Println("no pokemon in this area")
		return nil
	}

	for _, encounter := range result.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandHelp(context *Context, args []string) error {
	_, err := fmt.Println("Welcome to the Pokedex!\nUsage:")
	if err != nil {
		return err
	}
	for _, val := range COMMANDS {
		fmt.Printf("  %s: %s\n", val.name, val.description)
	}
	return nil
}

func commandExit(context *Context, args []string) error {
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
