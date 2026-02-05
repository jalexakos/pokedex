package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	pokecache "github.com/jalexakos/pokedex/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(map[string]cliCommand, *config, *pokecache.Cache, string, map[string]pokemon) error
}

type config struct {
	Next     string
	Previous string
}

type locationAreas struct {
	Count    int            `json:"count"`
	Next     string         `json:"next,omitempty"`
	Previous string         `json:"previous,omitempty"`
	Areas    []locationArea `json:"results"`
}

type locationArea struct {
	Name string
	URL  string
}

type exploreLocationAreas struct {
	Id                int                `json:"id"`
	Name              string             `json:"name"`
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

type pokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
	VersionDetails []versionDetail `json:"version_details"`
}

type versionDetail struct {
	Version struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"version"`
	Rate int `json:"rate"`
}

type pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []stat        `json:"stats"`
	Types          []pokemontype `json:"types"`
}

type stat struct {
	PokeStat pokestat `json:"stat"`
	BaseStat int      `json:"base_stat"`
}

type pokestat struct {
	Name string `json:"name"`
}

type pokemontype struct {
	PokeType poketype `json:"type"`
}

type poketype struct {
	Name string `json:"name"`
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Displays a map of the Pokedex",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Displays the previous map of the Pokedex",
		callback:    commandMapBack,
	},
	"explore": {
		name:        "explore",
		description: "Displays a list of all the Pokemon in a certain area, based on the location area name or id",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Attempts to catch a Pokemon",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Displays detailed information about a Pokemon",
		callback:    commandInspect,
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		Next:     "",
		Previous: "",
	}
	cache := pokecache.NewCache(time.Second * 5)
	pokedex := make(map[string]pokemon)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			// handle EOF or error, maybe break the loop
			break
		}
		inputText := scanner.Text()
		words := cleanInput(inputText)
		firstWord := words[0]
		secondWord := ""
		if len(words) > 1 {
			secondWord = words[1]
		}
		if command, exists := commands[firstWord]; exists {
			command.callback(commands, cfg, cache, secondWord, pokedex)
		} else {
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func commandExit(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string, pokedex map[string]pokemon) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string, pokedex map[string]pokemon) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, info := range cmds {
		fmt.Printf("%s: %s\n", info.name, info.description)
	}
	return nil
}

func commandMap(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string, pokedex map[string]pokemon) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.Next != "" {
		url = cfg.Next
	}
	var body []byte
	body, ok := cache.Get(url)
	if !ok {
		var err error
		body, err = getCall(url)
		if err != nil {
			return err
		}
		cache.Add(url, body)
	}
	var locationAreas locationAreas
	if err := json.Unmarshal(body, &locationAreas); err != nil {
		return err
	}
	if locationAreas.Next != "" {
		cfg.Next = locationAreas.Next
	}
	if locationAreas.Previous != "" {
		cfg.Previous = locationAreas.Previous
	}
	for _, locationArea := range locationAreas.Areas {
		fmt.Printf("%s\n", locationArea.Name)
	}
	return nil
}

func commandMapBack(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string, pokedex map[string]pokemon) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.Previous != "" {
		url = cfg.Previous
	}
	var body []byte
	fmt.Println("url:", url)
	body, ok := cache.Get(url)
	if !ok {
		var err error
		body, err = getCall(url)
		if err != nil {
			return err
		}
		cache.Add(url, body)
	}
	var locationAreas locationAreas
	if err := json.Unmarshal(body, &locationAreas); err != nil {
		return err
	}
	if locationAreas.Next != "" {
		cfg.Next = locationAreas.Next
	}
	if locationAreas.Previous != "" {
		cfg.Previous = locationAreas.Previous
	}
	for _, locationArea := range locationAreas.Areas {
		fmt.Printf("%s\n", locationArea.Name)
	}
	return nil
}

func getCall(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func commandExplore(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string, pokedex map[string]pokemon) error {
	url := "https://pokeapi.co/api/v2/location-area/" + secondWord
	var body []byte
	body, ok := cache.Get(url)
	if !ok {
		var err error
		body, err = getCall(url)
		if err != nil {
			return err
		}
		cache.Add(url, body)
	}
	var exploreLocationAreas exploreLocationAreas
	if err := json.Unmarshal(body, &exploreLocationAreas); err != nil {
		return err
	}

	exploringLine := "Exploring " + secondWord + "..."
	fmt.Println(exploringLine)
	for _, pokemonEncounter := range exploreLocationAreas.PokemonEncounters {
		fmt.Printf("%s\n", pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string, pokedex map[string]pokemon) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + secondWord
	var body []byte
	body, ok := cache.Get(url)
	if !ok {
		var err error
		body, err = getCall(url)
		if err != nil {
			return err
		}
		cache.Add(url, body)
	}
	var pokemon pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return err
	}
	catchingLine := "Throwing a Pokeball at " + pokemon.Name + "..."
	fmt.Println(catchingLine)
	attempt := rand.Intn(346)
	if attempt > pokemon.BaseExperience {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func commandInspect(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string, pokedex map[string]pokemon) error {
	if pokemon, ok := pokedex[secondWord]; ok {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%s: %d\n", stat.PokeStat.Name, stat.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, poketype := range pokemon.Types {
			fmt.Printf("  -%s\n", poketype.PokeType.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}
