package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(map[string]cliCommand) error
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
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			// handle EOF or error, maybe break the loop
			break
		}
		inputText := scanner.Text()
		firstWord := cleanInput(inputText)[0]
		if command, exists := commands[firstWord]; exists {
			command.callback(commands)
		} else {
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func commandExit(cmds map[string]cliCommand) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cmds map[string]cliCommand) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, info := range cmds {
		fmt.Printf("%s: %s\n", info.name, info.description)
	}
	return nil
}
