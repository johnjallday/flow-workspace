package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/project/todo"
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
			// Launch the TODO REPL for this project's todo.md
			todoFilePath := filepath.Join(projectDir, "todo.md")
			todo.StartTodoREPL(todoFilePath)

		case "launch":
			// Placeholder for launching the project (IDE, Docker, etc.)
			LaunchProject(projectDir, "Project")

		case "tree":
			// Print directory tree for the project using PrintTree from tree.go
			fmt.Printf("Directory tree for: %s\n", projectDir)
			PrintTree(projectDir, "")

		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

// printProjectHelp displays the available commands in this project-level REPL.
func printProjectHelp() {
	fmt.Println(`Available commands (Project REPL):
  todo    - Open the TODO REPL for this project's todo.md
  launch  - Launch the project
  tree    - Print directory tree of the project
  help    - Show this help message
  exit    - Exit the Project REPL
`)
}
