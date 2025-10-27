package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort" // --- Hint: Import the 'sort' package ---
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

func checkArgs() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %sgo%s run %s./src/%s <json-file-path>", Bold, Reset, Bold, Reset)
	}
}

func main() {
	checkArgs()

	// 1. Read the JSON file from arguments
	// --- Hint: Always check for errors! ---
	// os.ReadFile returns the data AND an error.
	jsonData, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Error reading file %s: %v", os.Args[1], err)
	}

	// 2. Define a variable to hold the unmarshalled data.
	// The structure is a map where the key is a string (the item name)
	// and the value is our ItemDetails struct.
	var items map[string]ItemDetails

	// 3. Unmarshal (parse) the JSON data.
	// We pass the JSON data as a byte slice and a pointer to our 'items' map.
	err = json.Unmarshal([]byte(jsonData), &items)
	if err != nil {
		// If the JSON is invalid, the program will stop here.
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// 4. --- NEW: Get and sort the keys for a stable order ---
	// Create a slice to hold all the keys from the map
	var itemNames []string
	// Loop over the map just to get the keys
	for name := range items {
		itemNames = append(itemNames, name)
	}
	// Now, sort the slice of keys alphabetically
	sort.Strings(itemNames)

	// 5. Iterate through the *sorted slice* and display the list.
	fmt.Println("--- Your Todo List ---")

	// We now loop over our sorted 'itemNames' slice, not the 'items' map
	for _, itemName := range itemNames {
		// --- Hint: Get the details from the map using the sorted key ---
		details := items[itemName]

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


