package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	pokecache "github.com/jalexakos/pokedex/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(map[string]cliCommand, *config, *pokecache.Cache, string) error
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
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		Next:     "",
		Previous: "",
	}
	cache := pokecache.NewCache(time.Second * 5)
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
			command.callback(commands, cfg, cache, secondWord)
		} else {
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func commandExit(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, info := range cmds {
		fmt.Printf("%s: %s\n", info.name, info.description)
	}
	return nil
}

func commandMap(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string) error {
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

func commandMapBack(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string) error {
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

func commandExplore(cmds map[string]cliCommand, cfg *config, cache *pokecache.Cache, secondWord string) error {
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

	fmt.Println("Exploring ", secondWord, "...")
	for _, pokemonEncounter := range exploreLocationAreas.PokemonEncounters {
		fmt.Printf("%s\n", pokemonEncounter.Pokemon.Name)
	}
	return nil
}
