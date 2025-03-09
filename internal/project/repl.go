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

		// Load and print project info.
		metaFile := filepath.Join(projectDir, "project_info.toml")
		proj, err := LoadProjectInfo(metaFile)
		if err != nil {
			fmt.Printf("Error loading project info: %v\n", err)
		} else {
			printProjectInfo(proj)
		}

		// Load and print todos.
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
		case "edit":
			// Use the project metadata file.
			if err := editProjectInfo(metaFile); err != nil {
				fmt.Printf("Error editing project info: %v\n", err)
			}
			fmt.Println("Press Enter to continue...")
			reader.ReadString('\n')
		case "todo":
			todo.StartTodoREPL(dbPath, todoFile)
			return
		case "exit":
			fmt.Println("Exiting Project REPL. Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
			fmt.Println("Press Enter to continue...")
			reader.ReadString('\n')
		}
	}
}

// printProjectInfo displays key project metadata on the screen.
func printProjectInfo(proj *Project) {
	fmt.Println("====================================")
	fmt.Println("Project Info:")
	fmt.Println("Name:         ", proj.Name)
	fmt.Println("Alias:        ", proj.Alias)
	fmt.Println("Project Type: ", proj.ProjectType)
	fmt.Println("Tags:         ", strings.Join(proj.Tags, ", "))
	fmt.Println("Notes:        ", strings.Join(proj.Notes, ", "))
	fmt.Printf("Date Created:  %v\n", proj.DateCreated)
	fmt.Printf("Date Modified: %v\n", proj.DateModified)
	if proj.GitURL != "" {
		fmt.Println("Git URL:      ", proj.GitURL)
	}
	fmt.Println("====================================")
}

// The rest of your functions remain unchanged...
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
  edit      - Edit project info (tags, notes, name, alias, project type)
  exit      - Exit the Project REPL`)
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("CLS for", runtime.GOOS, "not implemented")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
