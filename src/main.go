package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ItemDetails holds the properties for each item.
// We use struct tags (`json:"..."`) to map the
// JSON keys (like "check") to our Go struct fields (like "Check").
type ItemDetails struct {
	Check     bool `json:"check"`
	Important bool `json:"important"`
}

// ANSI escape codes for text formatting
// --- Hint: These are special strings your terminal understands. ---
// \x1b[1m tells the terminal "start bold text".
// \x1b[0m tells the terminal "reset text to normal".
const (
	Bold  = "\x1b[1m"
	Reset = "\x1b[0m"
)

func main() {
	// 1. Define the JSON data as a string.
	// I've expanded on your example with a few more items
	// to properly test the logic.
	jsonData := `
	{
	  "Buy Groceries": {
		"check": true,
		"important": true
	  },
	  "Call Mom": {
		"check": true,
		"important": false
	  },
	  "Pay Rent": {
		"check": true,
		"important": true
	  },
	  "Walk the Dog": {
		"check": false,
		"important": true
	  },
	  "Finish Report": {
		"check": false,
		"important": false
	  }
	}
	`

	// 2. Define a variable to hold the unmarshalled data.
	// The structure is a map where the key is a string (the item name)
	// and the value is our ItemDetails struct.
	var items map[string]ItemDetails

	// 3. Unmarshal (parse) the JSON data.
	// We pass the JSON data as a byte slice and a pointer to our 'items' map.
	err := json.Unmarshal([]byte(jsonData), &items)
	if err != nil {
		// If the JSON is invalid, the program will stop here.
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// 4. Iterate through the map and display it like a CLI list.
	fmt.Println("--- Your Todo List ---")

	// We'll iterate through all items now, not just the filtered ones.
	for itemName, details := range items {
		// --- Hint 1: Determine the checkbox prefix ---
		// We use a simple 'if' statement to decide what string to show.
		checkboxPrefix := "[ ]" // Default to unchecked
		if details.Check {
			checkboxPrefix = "[x]" // Set to checked if the 'check' bool is true
		}

		// --- Hint 2: Build the final output string ---
		// We check if the 'important' bool is true.
		if details.Important {
			// If it's important, we wrap the itemName with our ANSI codes.
			// fmt.Printf is used for formatted printing.
			// %s is a placeholder for a string. We pass variables for each %s.
			// The \n at the end means "print a new line".
			fmt.Printf("%s %s%s%s\n", checkboxPrefix, Bold, itemName, Reset)
		} else {
			// If not important, just print the prefix and the name.
			fmt.Printf("%s %s\n", checkboxPrefix, itemName)
		}
	}
}

