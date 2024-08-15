package models

import (
	"ass1/auth"
	"ass1/hash"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	state       state
	username    string
	password    string
	clearance   int
	loggedIn    bool
	fileSystem  map[string]File
	currentFile string
	input       string // Separate input field
	errMsg      string
	successMsg  string
}

type File struct {
	owner          string
	classification int
	content        string
}

type state int

const (
	initialState state = iota
	usernameState
	passwordState
	menuState
	createState
	appendState
	readState
	writeState
	listState
	saveState
	exitState
)

// InitializeUser handles the creation of a new user
func InitializeUser() error {
	var username, password, confirmPassword string
	var clearance int

	fmt.Print("Username: ")
	fmt.Scanln(&username)

	// Check if username already exists
	if userExists(username) {
		return errors.New("username already exists")
	}

	fmt.Print("Password: ")
	fmt.Scanln(&password)

	fmt.Print("Confirm Password: ")
	fmt.Scanln(&confirmPassword)

	if password != confirmPassword {
		return errors.New("passwords do not match")
	}

	// Check password requirements
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	fmt.Print("User clearance (0, 1, 2, 3): ")
	fmt.Scanln(&clearance)

	if clearance < 0 || clearance > 3 {
		return errors.New("invalid clearance level")
	}

	// Generate salt and hash password
	rand.Seed(time.Now().UnixNano())
	salt := fmt.Sprintf("%08d", rand.Intn(100000000))
	passSaltHash := hash.MD5Hash(password + salt)

	// Save to salt.txt and shadow.txt
	if err := saveToFile("salt.txt", fmt.Sprintf("%s:%s\n", username, salt)); err != nil {
		return err
	}
	if err := saveToFile("shadow.txt", fmt.Sprintf("%s:%s:%d\n", username, passSaltHash, clearance)); err != nil {
		return err
	}

	fmt.Println("User created successfully!")
	return nil
}

func userExists(username string) bool {
	file, err := os.ReadFile("salt.txt")
	if err != nil {
		return false
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, username+":") {
			return true
		}
	}
	return false
}

func saveToFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

// InitialModel returns the initial state of the model
func InitialModel() model {
	return model{
		state:      initialState,
		fileSystem: make(map[string]File),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": // Handle Ctrl+C signal
			return m, tea.Quit
		}

		switch m.state {
		case initialState:
			m.state = usernameState
			m.input = ""
			return m, nil

		case usernameState:
			m.handleTextInput(msg, &m.username, passwordState)
			return m, nil

		case passwordState:
			m.handlePasswordInput(msg)
			return m, nil

		case menuState:
			return m.handleMenu(msg)

		case createState:
			m.handleCreate(msg)
			return m, nil

		case appendState:
			m.handleAppend(msg)
			return m, nil

		case readState:
			m.handleRead(msg)
			return m, nil

		case writeState:
			m.handleWrite(msg)
			return m, nil

		case listState:
			m.input = m.listFiles()
			m.state = menuState
			return m, nil

		case saveState:
			err := m.saveFileSystem()
			if err != nil {
				m.errMsg = "Failed to save files: " + err.Error()
			} else {
				m.successMsg = "Files saved."
			}
			m.state = menuState
			return m, nil

		case exitState:
			if msg.String() == "y" || msg.String() == "Y" {
				return m, tea.Quit
			}
			m.state = menuState
			return m, nil
		}

	case tea.WindowSizeMsg:
		// Optional: Handle window resizing if needed
	}

	return m, nil
}

func (m model) View() string {
	output := ""

	if m.errMsg != "" {
		output += "Error: " + m.errMsg + "\n"
		m.errMsg = ""
	}

	if m.successMsg != "" {
		output += "Success: " + m.successMsg + "\n"
		m.successMsg = ""
	}

	switch m.state {
	case initialState:
		return output + "Welcome to the Secure File System\nPress any key to continue...\n"

	case usernameState:
		return output + "Enter Username: " + m.input

	case passwordState:
		return output + "Enter Password: " + strings.Repeat("*", len(m.input))

	case menuState:
		return output + "Options: (C)reate, (A)ppend, (R)ead, (W)rite, (L)ist, (S)ave or (E)xit.\n"

	case createState:
		return output + "Enter Filename to Create: " + m.input

	case appendState:
		return output + "Enter Filename to Append: " + m.input

	case readState:
		return output + "Enter Filename to Read: " + m.input

	case writeState:
		return output + "Enter Filename to Write: " + m.input

	case listState:
		return output + m.input

	case saveState:
		return output + "Saving files...\n"

	case exitState:
		return output + "Shut down the FileSystem? (Y)es or (N)o\n"

	default:
		return output + "Unknown state\n"
	}
}

func (m *model) handleMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "c":
		m.state = createState
	case "a":
		m.state = appendState
	case "r":
		m.state = readState
	case "w":
		m.state = writeState
	case "l":
		m.state = listState
	case "s":
		m.state = saveState
	case "e":
		m.state = exitState
	default:
		m.errMsg = "Invalid option. Try again."
	}
	return m, nil
}

func (m *model) handleCreate(msg tea.KeyMsg) {
	if msg.Type == tea.KeyEnter {
		filename := strings.TrimSpace(m.input) // Trim spaces to avoid issues with leading/trailing spaces
		if filename == "" {
			m.errMsg = "Filename cannot be empty."
		} else if _, exists := m.fileSystem[filename]; exists {
			m.errMsg = "File already exists."
		} else {
			m.fileSystem[filename] = File{owner: m.username, classification: m.clearance, content: ""}
			m.successMsg = fmt.Sprintf("File '%s' created successfully.", filename) // Clear and specific message
		}
		m.input = ""
		m.state = menuState
	} else if msg.Type == tea.KeyBackspace {
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	} else {
		m.input += msg.String()
	}
}

func (m *model) handleAppend(msg tea.KeyMsg) {
	if msg.Type == tea.KeyEnter {
		filename := m.input
		file, exists := m.fileSystem[filename]
		if !exists {
			m.errMsg = "File does not exist."
		} else if file.classification >= m.clearance {
			m.errMsg = "Access denied."
		} else {
			file.content += "\n" + m.input
			m.fileSystem[filename] = file
			m.successMsg = "Content appended."
		}
		m.input = ""
		m.state = menuState
	} else if msg.Type == tea.KeyBackspace {
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	} else {
		m.input += msg.String()
	}
}

func (m *model) handleRead(msg tea.KeyMsg) {
	if msg.Type == tea.KeyEnter {
		filename := m.input
		file, exists := m.fileSystem[filename]
		if !exists {
			m.errMsg = "File does not exist."
		} else if file.classification > m.clearance {
			m.errMsg = "Access denied."
		} else {
			m.successMsg = "File content: " + file.content
		}
		m.input = ""
		m.state = menuState
	} else if msg.Type == tea.KeyBackspace {
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	} else {
		m.input += msg.String()
	}
}

func (m *model) handleWrite(msg tea.KeyMsg) {
	if msg.Type == tea.KeyEnter {
		filename := m.input
		file, exists := m.fileSystem[filename]
		if !exists {
			m.errMsg = "File does not exist."
		} else if file.classification > m.clearance {
			m.errMsg = "Access denied."
		} else {
			file.content = m.input
			m.fileSystem[filename] = file
			m.successMsg = "File content overwritten."
		}
		m.input = ""
		m.state = menuState
	} else if msg.Type == tea.KeyBackspace {
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	} else {
		m.input += msg.String()
	}
}

func (m *model) listFiles() string {
	var builder strings.Builder
	for filename, file := range m.fileSystem {
		if file.classification <= m.clearance { // Enforce "no read up" policy
			builder.WriteString(fmt.Sprintf("File: %s, Owner: %s, Classification: %d\n", filename, file.owner, file.classification))
		}
	}
	if builder.Len() == 0 {
		return "No accessible files found.\n"
	}
	return builder.String()
}

func (m *model) saveFileSystem() error {
	f, err := os.Create("Files.store")
	if err != nil {
		return err
	}
	defer f.Close()

	for filename, file := range m.fileSystem {
		// Ensure filenames and file contents are properly formatted
		_, err := fmt.Fprintf(f, "%s:%s:%d:%s\n", strings.TrimSpace(filename), file.owner, file.classification, file.content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *model) LoadFileSystem() error {
	file, err := os.ReadFile("Files.store")
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			parts := strings.SplitN(line, ":", 4)
			if len(parts) < 4 {
				continue // Skip invalid entries
			}
			classification, err := strconv.Atoi(parts[2])
			if err != nil || classification < 0 || classification > 3 {
				continue // Skip invalid entries
			}
			m.fileSystem[strings.TrimSpace(parts[0])] = File{
				owner:          parts[1],
				classification: classification,
				content:        parts[3],
			}
		}
	}
	return nil
}

func (m *model) handleTextInput(msg tea.KeyMsg, target *string, nextState state) {
	switch msg.Type {
	case tea.KeyBackspace:
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	case tea.KeyEnter:
		*target = m.input
		m.input = ""
		m.state = nextState
	default:
		m.input += msg.String()
	}
}

func (m *model) handlePasswordInput(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyBackspace:
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	case tea.KeyEnter:
		m.password = m.input
		m.input = ""
		authSuccess, clearance, err := auth.AuthenticateUser(m.username, m.password)
		if err != nil {
			m.errMsg = "Authentication failed: " + err.Error()
			m.state = initialState
		} else if authSuccess {
			m.clearance = clearance
			m.loggedIn = true
			m.state = menuState
		} else {
			m.errMsg = "Invalid credentials"
			m.state = initialState
		}
	default:
		m.input += msg.String()
	}
}
