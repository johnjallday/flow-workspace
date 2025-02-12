package agent

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Agent represents the configuration for the CLI agent.
// Extend this struct with any additional configuration fields you need.
type Agent struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	// Add more fields as required.
}

// Command represents a CLI command.
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// getCommands executes the binary with the "--list-commands" flag and returns the list of commands.
func getCommands(agentPath string) ([]Command, error) {
	cmd := exec.Command(agentPath, "--list-commands")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing binary: %w", err)
	}

	var commands []Command
	if err := json.Unmarshal(output, &commands); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}
	return commands, nil
}

// printCommands prints the available commands in a style similar to your project REPL.
func printCommands(commands []Command) {
	fmt.Println("Available commands (Agent REPL):")

	// Sort the commands alphabetically.
	sortedCommands := make([]Command, len(commands))
	copy(sortedCommands, commands)
	sort.Slice(sortedCommands, func(i, j int) bool {
		return sortedCommands[i].Name < sortedCommands[j].Name
	})

	for _, cmd := range sortedCommands {
		fmt.Printf("  %-10s - %s\n", cmd.Name, cmd.Description)
	}
}

// LaunchAgent executes the agent binary with the specified command.
// It pipes the output to stdout/stderr.
func LaunchAgent(agentPath string, commandName string) {
	cmd := exec.Command(agentPath, commandName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing command %s: %v\n", commandName, err)
	}
}

// StartAgentREPL starts an interactive REPL for the agent binary.
// It fetches the available commands from the agent binary and waits for user input.
func StartAgentREPL(agentPath string) {
	// Ensure agentPath is absolute.
	absAgentPath, err := filepath.Abs(agentPath)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		return
	}

	commandsList, err := getCommands(absAgentPath)
	if err != nil {
		fmt.Println("Error retrieving commands:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Agent REPL started for binary: %s\n", absAgentPath)
	printCommands(commandsList)

	for {
		fmt.Printf("\n[agent:%s] >> ", filepath.Base(absAgentPath))
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)

		switch strings.ToLower(line) {
		case "exit":
			fmt.Println("Exiting Agent REPL. Goodbye!")
			return
		case "help", "commands":
			printCommands(commandsList)
		default:
			// Check if the entered command exists in the list.
			var found bool
			for _, cmd := range commandsList {
				if strings.EqualFold(cmd.Name, line) {
					found = true
					fmt.Printf("Executing command: %s\n", cmd.Name)
					LaunchAgent(absAgentPath, cmd.Name)
					break
				}
			}
			if !found {
				fmt.Println("Unknown command. Type 'commands' for available commands.")
			}
		}
	}
}
