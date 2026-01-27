package main

import (
	"bufio"
	"fmt"
	"os"
)

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
		fmt.Printf("Your command was: %v\n", firstWord)
	}
}
