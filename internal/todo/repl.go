package todo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	db "github.com/johnjallday/flow-workspace/internal/db/todo"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// StartTodoREPL is the interactive REPL for a single todo.md file using TodoService.
func StartTodoREPL(dbPath string, todoFilePath string) {
	reader := bufio.NewReader(os.Stdin)

	// Initialize the database.
	mydb, err := db.InitDB(dbPath)
	if err != nil {
		fmt.Println("Error connecting to db:", err)
		return
	}

	// Migrate finished todos from the file to the database.
	MigrateFinishedTodos(todoFilePath, mydb)

	// Create an instance of TodoService.
	service := NewFileTodoService(todoFilePath)

	// REPL loop.
	for {
		clearScreen()
		fmt.Println("dbPath:", dbPath)
		printHelp()

		// List current todos.
		todos, err := service.ListTodos()
		if err != nil {
			fmt.Printf("Error loading todos: %v\n", err)
		} else {
			PrintTodos(todos)
		}

		fmt.Printf("\n[todo:%s] >> ", filepath.Base(todoFilePath))
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		command := strings.ToLower(parts[0])

		switch command {
		case "exit":
			fmt.Println("Exiting TODO REPL. Goodbye!")
			return
		case "add":
			// Prompt for description and due date.
			fmt.Print("Enter task description: ")
			description, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading description:", err)
				continue
			}
			description = strings.TrimSpace(description)
			if description == "" {
				fmt.Println("Task description cannot be empty.")
				continue
			}

			fmt.Print("Enter due date (YYYY-MM-DD) or leave empty: ")
			dueDate, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading due date:", err)
				continue
			}
			dueDate = strings.TrimSpace(dueDate)

			if err := service.AddTodo(description, dueDate); err != nil {
				fmt.Println("Error adding task:", err)
			} else {
				fmt.Println("Task added successfully.")
			}
		case "complete":
			if len(todos) == 0 {
				fmt.Println("No tasks available to complete.")
				break
			}
			fmt.Print("Enter the task number to complete: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			input = strings.TrimSpace(input)
			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > len(todos) {
				fmt.Println("Invalid task number.")
				continue
			}
			if err := service.CompleteTodo(index - 1); err != nil {
				fmt.Println("Error completing task:", err)
			} else {
				fmt.Println("Task marked as completed.")
			}
		case "delete":
			if len(todos) == 0 {
				fmt.Println("No tasks available to delete.")
				break
			}
			fmt.Print("Enter the task number to delete: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			input = strings.TrimSpace(input)
			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > len(todos) {
				fmt.Println("Invalid task number.")
				continue
			}
			if err := service.DeleteTodo(index - 1); err != nil {
				fmt.Println("Error deleting task:", err)
			} else {
				fmt.Println("Task deleted successfully.")
			}
		case "edit":
			if len(todos) == 0 {
				fmt.Println("No tasks available to edit.")
				break
			}
			fmt.Print("Enter the task number to edit: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			input = strings.TrimSpace(input)
			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > len(todos) {
				fmt.Println("Invalid task number.")
				continue
			}
			fmt.Print("Enter new description (leave empty to keep current): ")
			newDescription, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading description:", err)
				continue
			}
			newDescription = strings.TrimSpace(newDescription)
			fmt.Print("Enter new due date (YYYY-MM-DD, leave empty to keep current): ")
			newDueDate, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading due date:", err)
				continue
			}
			newDueDate = strings.TrimSpace(newDueDate)
			fmt.Print("Enter new status (ongoing/complete, leave empty to keep current): ")
			newStatus, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading status:", err)
				continue
			}
			newStatus = strings.TrimSpace(newStatus)
			if err := service.EditTodo(index-1, newDescription, newDueDate, newStatus); err != nil {
				fmt.Println("Error editing task:", err)
			} else {
				fmt.Println("Task edited successfully.")
			}
		case "weekly":
			fmt.Println("Running weekly review...")
			ReviewWeekly(todos)
			// (Add any weekly review functionality here.)
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}

		fmt.Print("Press Enter to continue...")
		_, _ = reader.ReadString('\n')
	}
}

func printHelp() {
	fmt.Println(`Available commands (TODO REPL):
  add       - Add a new task
  complete  - Mark a task as completed
  delete    - Delete a task
  edit      - Edit a task (update description, due date, and/or status)
  weekly    - Run the weekly review for this TODO file
  exit      - Exit the TODO REPL`)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
