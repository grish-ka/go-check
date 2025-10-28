package main

import (
	"encoding/json" // TUTORIAL: Added this import for creating the new file
	"fmt"
	"log"
	"os"
	"sort"
	"flag"

	tea "github.com/charmbracelet/bubbletea"
	// We will add the textinput import in a later step!
	ti "github.com/charmbracelet/bubbles/textinput" // <-- ADD THIS

	lg "github.com/charmbracelet/lipgloss"
)

var (
	// Green for checked items
	styleChecked = lg.NewStyle().Foreground(lg.Color("42")) // "42" is a nice green

	// Red/Bold for important, unchecked items
	styleImportant = lg.NewStyle().Foreground(lg.Color("203")).Bold(true) // "203" is a hot pink/red

	// Faint yellow for normal, unchecked items
	styleNormal = lg.NewStyle().Foreground(lg.Color("226")) // "226" is a bright yellow

	// Style for the cursor
	styleCursor = lg.NewStyle().Foreground(lg.Color("212")) // A bright magenta

	// Style for the help text
	styleHelp = lg.NewStyle().Foreground(lg.Color("240")) // A dim grey

	// --- TUTORIAL (STEP 6): Add a style for the popup box ---
	stylePopupBox = lg.NewStyle().
		Border(lg.RoundedBorder()).
		BorderForeground(lg.Color("63")). // A nice purple
		Padding(1)

)

// TUTORIAL (STEP 4): Define our app's "states"
type appState int

const (
    stateListBrowse appState = iota // 0
    stateItemCreate                 // 1
)

// --- Step 1: Define your Model ---
// (This is your original model, unchanged)
type model struct {
	items     map[string]ItemDetails // Your existing map
	itemNames []string               // Your sorted list of keys
	cursor    int                    // Which item we're pointing at
	filePath  string                 // The path to the file we're editing
	// --- TUTORIAL (STEP 4): Add new fields for "create" mode ---
    state            appState      // Tracks if we are browsing or creating
    textInput        ti.Model      // The textbox component
    newItemImportant bool          // Tracks the "important" toggle for the new item

	// --- TUTORIAL (STEP 6): Add window size fields ---
	width            int           // Holds the current terminal width
	height           int           // Holds the current terminal height
}

// ItemDetails holds the properties for each item.
type ItemDetails struct {
	Check     bool `json:"check"`
	Important bool `json:"important"`
}

// --- Define Command-Line Flags ---
var filePathPtr = flag.String("file", "", "Path to your JSON todo file (e.g., -file todos.json)")
var newFlagPtr = flag.String("new", "", "Create a new todo list with this as the first item (e.g., --new 'My task')")

// getFilePath finds the file path from flags or arguments
func getFilePath() string {
	// 1. Check if the -file flag was used
	if *filePathPtr != "" {
		return *filePathPtr
	}

	// 2. If not, check if there's a positional argument
	if flag.NArg() > 0 {
		return flag.Arg(0)
	}

	// 3. If no file is provided, show a helpful error and exit.
	log.Println("Error: Missing file argument.")
	fmt.Println("Usage (View): go run . <path-to-file.json>")
	fmt.Println("Usage (New):  go run . --new <'first item'> -file <path-to-file.json>")
	os.Exit(1) // Exit the program
	return ""  // This line will never be reached
}

// TUTORIAL: I removed your 'checkNewFile' function as we will
// put this logic directly in main() to make it clearer.

// --- Helper function to load the file ---
// (Unchanged)
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

// TUTORIAL (STEP 5): Add a save function
// This is a "method" on our model. It can access m.items and m.filePath.
func (m model) save() error {
    // 1. Convert our 'items' map back into pretty JSON
    //    We use MarshalIndent for a nice, readable file.
    jsonData, err := json.MarshalIndent(m.items, "", "    ")
    if err != nil {
        return err
    }

    // 2. Write the JSON data back to the original file
    //    0644 are standard file permissions.
    err = os.WriteFile(m.filePath, jsonData, 0644)
    if err != nil {
        return err
    }

    return nil
}

// --- Step 2: Create your initial model ---
// (Unchanged)
// --- Step 2: Create your initial model ---
func initialModel(filePath string) model {
    items, itemNames := loadData(filePath)

    // TUTORIAL (STEP 4): Create our new text input
    txtInput := ti.New()                 // Create a new textbox
    txtInput.Placeholder = "New todo..." // Set placeholder text
    txtInput.CharLimit = 200             // Set a max length

    return model{
        items:     items,
        itemNames: itemNames,
        cursor:    0,
        filePath:  filePath,

        // --- TUTORIAL (STEP 4): Set the initial state ---
        state:            stateListBrowse, // Start in "browse" mode
        textInput:        txtInput,
        newItemImportant: false,
    }
}

// --- Step 3: Define your Init function ---
// (Unchanged)
func (m model) Init() tea.Cmd {
	return nil
}

// --- Step 4: Define your Update function ---
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// --- TUTORIAL (STEP 6): Handle window size messages ---
	// This switch handles messages that we care about *regardless* of our app state.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// This is your *existing* state machine
	switch m.state {

    // --- STATE 1: BROWSING THE LIST ---
    // This is where all your existing keybinds go
    case stateListBrowse:
        switch msg := msg.(type) {
        case tea.KeyMsg:
            // Which key was pressed?
            switch msg.String() {

            // These are all your keybinds from Steps 2 & 3
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

				if err := m.save(); err != nil { log.Printf("Error saving: %v", err) }
            
            // TUTORIAL (STEP 2): Add the 'i' keybind
            case "i":
                itemName := m.itemNames[m.cursor]
                details := m.items[itemName]
                details.Important = !details.Important // <-- The only change is here
                m.items[itemName] = details

				if err := m.save(); err != nil { log.Printf("Error saving: %v", err) }
            
            // TUTORIAL (STEP 3): Add the 'r' keybind
            case "r":
                // Safety check: Don't do anything if the list is already empty
                if len(m.itemNames) == 0 {
                    return m, nil
                }

                // 1. Get the name of the item to delete
                itemName := m.itemNames[m.cursor]

                // 2. Delete it from the map (the easy part)
                delete(m.items, itemName)

                // 3. Delete it from the slice (the "slice trick")
                m.itemNames = append(m.itemNames[:m.cursor], m.itemNames[m.cursor+1:]...)

                // 4. Fix the cursor so it doesn't go "out of bounds"
                if m.cursor >= len(m.itemNames) {
                    m.cursor = len(m.itemNames) - 1
                }

                // If we just emptied the list, set cursor to 0
                if m.cursor < 0 {
                    m.cursor = 0
                }
				
				if err := m.save(); err != nil { log.Printf("Error saving: %v", err) }
            // TUTORIAL (STEP 4): 'n' keybind to switch modes
            case "n":
                m.state = stateItemCreate      // 1. Change to "create" mode
                m.newItemImportant = false     // 2. Reset the toggle
                m.textInput.Reset()            // 3. Clear the textbox
                m.textInput.Focus()            // 4. Focus the textbox
                return m, ti.Blink             // 5. The "trick": return a command to make the cursor blink!
            }
        }
        return m, nil

    // --- STATE 2: CREATING A NEW ITEM ---
    // This is all new logic for the popup
    case stateItemCreate:
        var cmd tea.Cmd // This will hold any commands from the textbox

        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {

            // "Enter" saves the new item
            case "enter":
                newItemName := m.textInput.Value()
                if newItemName != "" {
                    // Add the new item to our map
                    m.items[newItemName] = ItemDetails{
                        Check:     false,
                        Important: m.newItemImportant,
                    }
                    // Add the new name to our slice
                    m.itemNames = append(m.itemNames, newItemName)
                    
                    // The "Gotcha": Re-sort the slice!
                    sort.Strings(m.itemNames)

					if err := m.save(); err != nil { log.Printf("Error saving: %v", err) }
                }
                // Go back to browse mode
                m.state = stateListBrowse
                return m, nil

            // "Esc" cancels and goes back
            case "esc":
                m.state = stateListBrowse
                return m, nil

            // "Tab" toggles the "important" switch
            case "tab":
                m.newItemImportant = !m.newItemImportant
                return m, nil
            }
        }

        // The "trick": This is how you type!
        // We pass the message to the textbox's own Update function.
        m.textInput, cmd = m.textInput.Update(msg)
        return m, cmd
    }

    return m, nil
}

// --- Step 5: Define your View function ---
// (Unchanged for now)
func (m model) View() string {
	// 's' will be our final string
	s := "--- Your Todo List ---\n\n"

	// Iterate over our sorted list, just like before
	for i, itemName := range m.itemNames {
		// Get the details from the map
		details := m.items[itemName]

		// --- Handle the cursor ---
		// 'cursorStr' will be "> " if this is the selected line, or "  " if not.
		// --- TUTORIAL (STEP 2): Style the list ---

		// 1. Handle the cursor
		cursorStr := "  " // Default empty space
		if m.cursor == i {
			cursorStr = styleCursor.Render("> ") // Apply cursor style
		}

		// 2. Handle the checkbox
		checkbox := "[ ]"
		if details.Check {
			checkbox = "[x]"
		}

		// 3. Combine checkbox and item name into one line
		line := fmt.Sprintf("%s %s", checkbox, itemName)

		// 4. Choose the correct style for the line based on its state
		if details.Check {
			line = styleChecked.Render(line)
		} else if details.Important {
			line = styleImportant.Render(line)
		} else {
			line = styleNormal.Render(line)
		}

		// 5. Add the cursor and the styled line to our final string
		s += fmt.Sprintf("%s%s\n", cursorStr, line)
	}

	// Add a footer with help info
	helpText := "\nPress 'j/down' and 'k/up' to move.\n"
	helpText += "Press 'space' to toggle check.\n"
    helpText += "Press 'i' to toggle important.\n" // <-- Make sure you have this
    helpText += "Press 'r' to remove an item.\n"   // <-- Make sure you have this
    helpText += "Press 'n' to add a new item.\n"   // <-- ADD THIS LINE
    helpText += "Press 'q' or 'ctrl+c' to quit.\n"

	s += styleHelp.Render(helpText)

	// Return the final string
	// --- TUTORIAL (STEP 6): Draw the popup as a modal ---
	
	// We'll return the main list `s` (which we built above)
	mainView := s

	// If we're in "create" mode, we'll draw the popup *on top*
	if m.state == stateItemCreate {
		
		// 1. Build the "[x] Important" switch string
		checkbox := "[ ]"
		if m.newItemImportant {
			checkbox = styleChecked.Render("[x]") // Use our green style!
		}

		// 2. Build the *content* of the popup
		popupContent := fmt.Sprintf(
			"--- Add New Item ---\n%s\n\n%s Important (press 'tab' to toggle)\n\n%s",
			m.textInput.View(), // This draws the actual textbox!
			checkbox,
			styleHelp.Render("(press 'enter' to save, 'esc' to cancel)"),
		)

		// 3. Style the content in our popup box
		styledPopup := stylePopupBox.Render(popupContent)

		// 4. Use lipgloss.Place() to center the popup.
		// This is the magic! It creates a new string with the
		// popup placed in the middle of the screen.
		centeredPopup := lg.Place(
			m.width,  // Full screen width
			m.height, // Full screen height
			lg.Center, // Horizontal alignment
			lg.Center, // Vertical alignment
			styledPopup,
		)

		// 5. Return the main view *plus* the centered popup.
		// This draws the popup "on top" of the list.
		return mainView + centeredPopup
	}

	// If we're not in "create" mode, just return the main list
	return mainView
}

// --- Step 6: Define your new main function ---
func main() {
	// 1. Parse flags
	flag.Parse()

	// 2. TUTORIAL: Handle the --new flag
	// This logic runs *before* the TUI starts.
	if *newFlagPtr != "" {
		// Get the file path (this will exit if no path is given)
		filePath := getFilePath()

		// Create a new map with the one item from the flag
		newItemName := *newFlagPtr
		newMap := make(map[string]ItemDetails)
		newMap[newItemName] = ItemDetails{Check: false, Important: false}

		// Convert the map to pretty JSON
		jsonData, err := json.MarshalIndent(newMap, "", "    ")
		if err != nil {
			log.Fatalf("Error creating JSON: %v", err)
		}

		// Write the new file
		if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			log.Fatalf("Error writing new file: %v", err)
		}

		// Success!
		fmt.Printf("Successfully created new file at %s\n", filePath)

		// Exit the program. The TUI will not start.
		os.Exit(0)
	}

	// 3. Get file path (this only runs if --new was *not* used)
	filePath := getFilePath()
	// ----------------

	// Create our initial model
	m := initialModel(filePath)

	// Start the Bubble Tea program
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error starting program: %v", err)
	}
}

