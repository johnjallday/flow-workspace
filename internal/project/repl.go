package project

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/agent"
	"github.com/johnjallday/flow-workspace/internal/todo"
)

// StartProjectREPL starts an interactive REPL for a single project directory.
func StartProjectREPL(dbPath string, projectDir string) {
	clearScreen()

	todoFile := filepath.Join(projectDir, "todo.md")
	service := todo.NewFileTodoService(todoFile)
	todos, err := service.ListTodos()
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
	} else {
		todo.PrintTodos(todos)
	}

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
			clearScreen()
			printProjectHelp()
		case "todo":
			// Exit the current REPL and then start the TODO REPL.
			todoFile := filepath.Join(projectDir, "todo.md")
			todo.StartTodoREPL(dbPath, todoFile)
			return
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
  todo    - Open the TODO REPL for this project (exits Project REPL)
  launch  - Launch the project
  help    - Show this help message
  agent   - Start the Agent REPL
  exit    - Exit the Project REPL`)
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("clear") // Linux
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") // Windows
	case "darwin":
		cmd = exec.Command("clear") // macOS
	default:
		fmt.Println("CLS for", runtime.GOOS, "not implemented")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
