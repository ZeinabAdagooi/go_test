package main

import (
	"fmt"
	"sort"
)

func rearrangeString(s string) string {
	// Count the frequency of each character
	charCount := make(map[rune]int)
	for _, char := range s {
		charCount[char]++
	}

	// Create a slice of characters and sort them by frequency
	keys := make([]rune, 0, len(charCount))
	for key := range charCount {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return charCount[keys[i]] > charCount[keys[j]]
	})

	// Check if it's possible to rearrange the string
	if 2*charCount[keys[0]]-1 > len(s) {
		return "" // Not possible
	}

	// Rearrange the string
	result := make([]rune, len(s))
	index := 0
	for _, key := range keys {
		for charCount[key] > 0 {
			result[index] = key
			index += 2
			if index >= len(s) {
				index = 1
			}
			charCount[key]--
		}
	}

	return string(result)
}

func promptUser(message string) string {
	var userInput string
	fmt.Print(message)
	fmt.Scanln(&userInput)
	return userInput
}
func main() {

	userInput := promptUser("Enter some text: ")
	result := rearrangeString(userInput)
	fmt.Printf("Result: %s\n", result)
}
