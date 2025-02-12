package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/agent"
	"github.com/johnjallday/flow-workspace/internal/todo"
)

// StartProjectREPL starts an interactive REPL for a single project directory.
func StartProjectREPL(projectDir string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Project REPL started for directory: %s\n", projectDir)
	printProjectHelp()

	for {
		fmt.Printf("\n[project:%s] >> ", filepath.Base(projectDir))
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)
		switch strings.ToLower(line) {
		case "exit":
			fmt.Println("Exiting Project REPL. Goodbye!")
			return
		case "help":
			printProjectHelp()
		case "todo":
			// Launch the centralized TODO REPL.
			todoFile := filepath.Join(projectDir, "todo.md")
			todo.StartTodoREPL(todoFile)
		case "launch":
			LaunchProject(projectDir, "Project")
		case "agent":
			agent.StartAgentREPL("/Users/jj/Workspace/johnj-programming/coder/main")
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func printProjectHelp() {
	fmt.Println(`Available commands (Project REPL):
  todo    - Open the TODO REPL for this project
  launch  - Launch the project
  help    - Show this help message
  agent   - Start the Agent REPL
  exit    - Exit the Project REPL`)
}
