package main

import "strings"

func cleanInput(text string) []string {
	loweredStrings := strings.ToLower(text)
	cleanedStrings := strings.Fields(loweredStrings)
	return cleanedStrings
}
