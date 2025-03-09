package project

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/todo"
)

// StartProjectREPL starts an interactive REPL for a single project directory.
func StartProjectREPL(dbPath string, projectDir string) {
	coderPath := "/Users/jj/Workspace/johnj-programming/gorani-coder/main"
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()

		todoFile := filepath.Join(projectDir, "todo.md")

		service := todo.NewFileTodoService(todoFile)
		todos, err := service.ListTodos()
		if err != nil {
			fmt.Printf("Error loading todos: %v\n", err)
		} else {
			todo.PrintTodos(todos)
		}

		printProjectHelp()
		fmt.Println("Enter a command:")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)

		switch strings.ToLower(line) {
		case "implement":
			implementTodo(service, coderPath, reader)
		case "finish":
			finishTodo(service, coderPath, reader)
		case "exit":
			fmt.Println("Exiting Project REPL. Goodbye!")
			return
		case "todo":
			todoFile := filepath.Join(projectDir, "todo.md")
			todo.StartTodoREPL(dbPath, todoFile)
			return
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func implementTodo(service *todo.FileTodoService, coderPath string, reader *bufio.Reader) {
	executeTodoCommand(service, coderPath, reader, "ongoing", "implement", "create")
}

func finishTodo(service *todo.FileTodoService, coderPath string, reader *bufio.Reader) {
	executeTodoCommand(service, coderPath, reader, "complete", "implement", "merge")
}

func executeTodoCommand(service *todo.FileTodoService, coderPath string, reader *bufio.Reader, status string, command string, action string) {
	todos, err := service.ListTodos()
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		return
	}
	if len(todos) == 0 {
		fmt.Println("No todos available.")
		return
	}

	fmt.Print("Enter the number of the todo: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(todos) {
		fmt.Println("Invalid todo number.")
		return
	}

	if err := service.EditTodo(index-1, "", "", status); err != nil {
		fmt.Printf("Error updating todo: %v\n", err)
		return
	}

	fmt.Printf("Todo updated successfully to %s!\n", status)
	description := todos[index-1].Description
	formattedDescription := strings.ReplaceAll(description, " ", "-")

	cmd := exec.Command(coderPath, command, action, formattedDescription)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	} else {
		fmt.Printf("Command '%s %s' executed successfully!\n", command, action)
	}
}

func printProjectHelp() {
	fmt.Println(`Available commands (Project REPL):
  todo      - Open the TODO REPL for this project (exits Project REPL)
  launch    - Launch the project
  implement - Implement a todo
  finish    - Mark a todo as complete
  exit      - Exit the Project REPL`)
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "darwin":
		cmd = exec.Command("clear")
	default:
		fmt.Println("CLS for", runtime.GOOS, "not implemented")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
