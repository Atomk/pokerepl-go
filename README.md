# pokeapi-go

A simple REPL to play with Pokemon data. Uses [PokeAPI](https://pokeapi.co/) with caching.

## Usage

Start the REPL with `go run .`, then use the command `help` to see a list of all available commands.

Example output:
```
Pokedex > pokedex
Your Pokedex is empty. Go catch some Pokemon!

Pokedex > catch charizard
Catch odds: 8.78%
Throwing a pokeball at Charizard...
Throwing a pokeball at Charizard...
Throwing a pokeball at Charizard...
Charizard fleed!

Pokedex > catch caterpie
Catch odds: 54.00%
Throwing a pokeball at Caterpie...
Throwing a pokeball at Caterpie...
Caught!
New Pokemon! Adding data to the Pokedex

Pokedex > inspect caterpie
Name: Caterpie
Height: 3
Weight: 29
Stats:
  - hp: 45
  - attack: 30
  - defense: 35
  - special-attack: 20
  - special-defense: 20
  - speed: 45
Types:
  - bug

Pokedex > help
...
```

## Testing

```sh
go test .
go test ./internal
```
