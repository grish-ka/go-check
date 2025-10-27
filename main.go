package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"flag"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Step 1: Define your Model ---
// The model holds your application's state.
// We'll move all the important data from your 'main' function here.
type model struct {
	items     map[string]ItemDetails // Your existing map
	itemNames []string               // Your sorted list of keys
	cursor    int                    // Which item we're pointing at
	filePath  string                 // The path to the file we're editing
}

// ItemDetails holds the properties for each item.
// (This is the same struct as before)
type ItemDetails struct {
	Check     bool `json:"check"`
	Important bool `json:"important"`
}

const (
	Bold  = "\x1b[1m"
	Reset = "\x1b[0m"
)
var filePathPtr = flag.String("file", "", "Path to your JSON todo file (e.g., -file todos.json)")
func getFilePath() string {
	// 1. Check if the -file flag was used
	if *filePathPtr != "" {
		return *filePathPtr
	}

	// 2. If not, check if there's a positional argument
	// flag.NArg() returns the number of args *after* flags
	if flag.NArg() > 0 {
		// flag.Arg(0) returns the first positional arg
		return flag.Arg(0)
	}

	// 3. If no file is provided in any way, show a helpful error
	log.Println("Error: Missing file argument.")
	fmt.Println("Usage 1: go run . -file <path-to-your-file.json>")
	fmt.Println("Usage 2: go run . <path-to-your-file.json>")
	os.Exit(1) // Exit the program
	return ""  // This line will never be reached, but it satisfies the compiler
}

// --- Helper function to load the file ---
// We'll call this once, when the app starts.
func loadData(filePath string) (map[string]ItemDetails, []string) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var items map[string]ItemDetails
	err = json.Unmarshal(jsonData, &items)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	var itemNames []string
	for name := range items {
		itemNames = append(itemNames, name)
	}
	sort.Strings(itemNames)

	return items, itemNames
}

// --- Step 2: Create your initial model ---
// This function is called once when the program starts.
func initialModel(filePath string) model {
	items, itemNames := loadData(filePath)

	return model{
		items:     items,
		itemNames: itemNames,
		cursor:    0, // Start cursor at the first item
		filePath:  filePath,
	}
}

// --- Step 3: Define your Init function ---
// Init is for running "commands" when the app starts.
// For now, we don't need to do anything here.
func (m model) Init() tea.Cmd {
	return nil
}

// --- Step 4: Define your Update function ---
// This is your main event loop! It's called every time
// something happens (like a key press).
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Which key was pressed?
		switch msg.String() {

		// These keys will exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up.
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down.
		case "down", "j":
			if m.cursor < len(m.itemNames)-1 {
				m.cursor++
			}

		// The "spacebar" key will toggle the 'check' status.
		case " ":
			// 1. Get the name of the item at the cursor
			itemName := m.itemNames[m.cursor]
			// 2. Get the details
			details := m.items[itemName]
			// 3. Flip the boolean
			details.Check = !details.Check
			// 4. Put the new details back in the map
			m.items[itemName] = details
		}
	}

	// Return the updated model to Bubble Tea
	return m, nil
}

// --- Step 5: Define your View function ---
// This is your *new* display loop. It's called every time
// the model changes. It returns a SINGLE string for the
// entire screen.
func (m model) View() string {
	// 's' will be our final string
	s := "--- Your Todo List ---\n\n"

	// Iterate over our sorted list, just like before
	for i, itemName := range m.itemNames {
		// Get the details from the map
		details := m.items[itemName]

		// --- Handle the cursor ---
		// 'cursorStr' will be "> " if this is the selected line, or "  " if not.
		cursorStr := "  "
		if m.cursor == i {
			cursorStr = "> "
		}

		// --- Handle the checkbox ---
		checkbox := "[ ]"
		if details.Check {
			checkbox = "[x]"
		}

		// --- Handle 'Important' (bold) ---
		// We can still use our ANSI codes!
		itemNameStr := itemName
		if details.Important {
			// Note: Your 'Bold' and 'Reset' consts would be defined up top
			itemNameStr = fmt.Sprintf("\x1b[1m%s\x1b[0m", itemName)
		}

		// Add the line to our final string 's'
		s += fmt.Sprintf("%s%s %s\n", cursorStr, checkbox, itemNameStr)
	}

	// Add a footer with help info
	s += "\nPress 'j/down' and 'k/up' to move.\n"
	s += "Press 'space' to toggle check.\n"
	s += "Press 'q' or 'ctrl+c' to quit.\n"

	// Return the final string
	return s
}

// --- Step 6: Define your new main function ---
func main() {
	// We can re-use your old 'getFilePath' logic from json_reader.go!
	// (You would need to copy the 'getFilePath' and 'filePathPtr'
	// 'flag' logic from your other file into this one)
	// For this example, I'll hardcode it, but you should replace this.
	// ----------------
	// 1. Parse flags
	flag.Parse()
	// 2. Get file path
	filePath := getFilePath() // <-- Make sure to copy getFilePath and filePathPtr
	// ----------------

	// Create our initial model
	m := initialModel(filePath)

	// Start the Bubble Tea program
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting program: %v", err)
	}
}
