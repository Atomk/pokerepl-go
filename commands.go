package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	pokecache "github.com/atomk/pokedexcli/internal"
)

// Contains URLs used for pagination.
// If both previous and Next are null, no request was made to get location areas.
type Context struct {
	cache    *pokecache.Cache
	Pokedex  map[string]Pokemon
	Previous *string
	Next     *string
}

func NewContext(minutes uint) *Context {
	return &Context{
		cache:   pokecache.NewCache(time.Duration(minutes) * time.Minute),
		Pokedex: map[string]Pokemon{},
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

func commandCatch(context *Context, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("command `catch` requires exactly one argument")
	}

	pokemonName := strings.TrimSpace(args[0])
	if len(pokemonName) == 0 {
		return fmt.Errorf("provided argument is an empty string")
	}

	result, err := getPokemon(pokemonName, context.cache)
	if err != nil {
		return fmt.Errorf("could not get data about `%s`: %v", pokemonName, err)
	}

	/*
		Base experience:
		- caterpie: 39
		- magikarp: 40
		- dragonite, mew: 270
		- mewtwo, rayquaza: 306
		- arceus: 324
		min: 39, max: 324

		So `324 / baseExp`` is always >= 1
		I want 324 to mean 1% probability to capture with a throw.
		With this formula minimum catch probability is 324/39 ~= 8.3.
		That's too low, so I multiply by a constant.
	*/
	threshold := (324 / float64(result.BaseExperience)) * 7
	if threshold > 100 {
		threshold = 100
	}
	fmt.Println("base experience:", result.BaseExperience)
	fmt.Println("threshold:", threshold)
	catched := false
	for _ = range 5 {
		fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

		if rand.Intn(100) <= int(threshold) {
			fmt.Println("Catched!")
			if _, ok := context.Pokedex[pokemonName]; !ok {
				fmt.Println("New Pokemon! Adding data to the Pokedex")
				context.Pokedex[pokemonName] = Pokemon{}
			} else {
				fmt.Println("You already catched this Pokemon")
			}
			catched = true
			break
		}
	}

	if !catched {
		fmt.Printf("%s fleed!\n", pokemonName)
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
