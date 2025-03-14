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

	db "github.com/johnjallday/flow-workspace/internal/db/todo"
	"github.com/johnjallday/flow-workspace/internal/todo"
)

// StartProjectREPL starts an interactive REPL for a single project directory.
func StartProjectREPL(dbPath string, projectDir string) {
	coderPath := "/Users/jj/Workspace/johnj-programming/gorani-coder/main"
	reader := bufio.NewReader(os.Stdin)

	mydb, err := db.InitDB(dbPath)
	if err != nil {
		fmt.Println("Error connecting to db:", err)
		return
	}

	todoFile := filepath.Join(projectDir, "todo.md")
	todo.MigrateFinishedTodos(todoFile, mydb)

	for {
		clearScreen()

		// Determine the metadata file path.
		metaFile := filepath.Join(projectDir, "project_info.toml")
		var proj *Project
		proj, err := LoadProjectInfo(metaFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("project_info.toml not found.")
				fmt.Print("Would you like to import this directory? (y/n): ")
				answer, _ := reader.ReadString('\n')
				answer = strings.TrimSpace(strings.ToLower(answer))
				if answer == "y" || answer == "yes" {
					if err := ImportProject(projectDir); err != nil {
						fmt.Printf("Error importing project: %v\n", err)
					} else {
						// Try to load project info again.
						proj, err = LoadProjectInfo(metaFile)
						if err != nil {
							fmt.Printf("Error loading project info after import: %v\n", err)
						}
					}
				}
			} else {
				fmt.Printf("Error loading project info: %v\n", err)
			}
		}
		if proj != nil {
			printProjectInfo(proj)
		}

		// Load and print todos.
		service := todo.NewFileTodoService(todoFile)
		todos, err := service.ListTodos()
		if err != nil {
			fmt.Printf("Error loading todos: %v\n", err)
		} else {
			ongoingTodos := todo.FilterTodosByOngoing(todos)
			if len(ongoingTodos) == 0 {
				fmt.Println("No ongoing tasks found.")
				todo.PrintTodos(todos)
			} else {
				todo.PrintTodos(ongoingTodos)
			}
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
			if err := editProjectInfo(metaFile); err != nil {
				fmt.Printf("Error editing project info: %v\n", err)
			}
			fmt.Println("Press Enter to continue...")
			reader.ReadString('\n')
		case "todo":
			todo.StartTodoREPL(dbPath, todoFile)
			return
		case "add-todo":
			fmt.Print("Enter todo description: ")
			description, _ := reader.ReadString('\n')
			description = strings.TrimSpace(description)

			fmt.Print("Enter due date (YYYY-MM-DD) or leave blank: ")
			dueDate, _ := reader.ReadString('\n')
			dueDate = strings.TrimSpace(dueDate)

			// Reuse the same business logic
			err := AddTodoToProject(projectDir, description, dueDate)
			if err != nil {
				fmt.Printf("Error adding todo: %v\n", err)
			} else {
				fmt.Println("Todo added successfully!")
			}

			fmt.Println("Press Enter to continue...")
			reader.ReadString('\n')
		case "edit-todo":
			todos, err := service.ListTodos()
			if err != nil {
				fmt.Printf("Error loading todos: %v\n", err)
				break
			}

			if len(todos) == 0 {
				fmt.Println("No tasks available to edit.")
				break
			}

			// Prompt user to pick a todo by number
			fmt.Print("Enter the number of the todo to edit: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				break
			}
			input = strings.TrimSpace(input)
			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > len(todos) {
				fmt.Println("Invalid todo number.")
				break
			}
			selectedIndex := index - 1

			// Prompt for new values (leave blank to keep current)
			fmt.Printf("Current description: %s\n", todos[selectedIndex].Description)
			fmt.Print("Enter new description (leave blank to keep current): ")
			newDescription, _ := reader.ReadString('\n')
			newDescription = strings.TrimSpace(newDescription)

			fmt.Printf("Current due date: %s\n", todos[selectedIndex].DueDate.Format("2006-01-02"))
			fmt.Print("Enter new due date (YYYY-MM-DD, leave blank to keep current): ")
			newDueDate, _ := reader.ReadString('\n')
			newDueDate = strings.TrimSpace(newDueDate)

			fmt.Print("Enter new status (ongoing/complete, leave blank to keep current): ")
			newStatus, _ := reader.ReadString('\n')
			newStatus = strings.TrimSpace(newStatus)

			// Call the service to edit the todo
			err = service.EditTodo(selectedIndex, newDescription, newDueDate, newStatus)
			if err != nil {
				fmt.Printf("Error editing todo: %v\n", err)
			} else {
				fmt.Println("Todo edited successfully.")
			}

			fmt.Println("Press Enter to continue...")
			reader.ReadString('\n')

		case "delete-todo":
			todos, err := service.ListTodos()
			if err != nil {
				fmt.Printf("Error loading todos: %v\n", err)
				break
			}

			if len(todos) == 0 {
				fmt.Println("No tasks available to delete.")
				break
			}

			todo.PrintTodos(todos)

			fmt.Print("Enter the number of the todo to delete: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				break
			}
			input = strings.TrimSpace(input)
			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > len(todos) {
				fmt.Println("Invalid task number.")
				break
			}
			selectedIndex := index - 1

			// Confirm deletion
			fmt.Printf("Are you sure you want to delete task #%d: \"%s\"? (y/n): ", index, todos[selectedIndex].Description)
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(strings.ToLower(confirm))
			if confirm != "y" && confirm != "yes" {
				fmt.Println("Delete canceled.")
				break
			}

			if err := service.DeleteTodo(selectedIndex); err != nil {
				fmt.Printf("Error deleting task: %v\n", err)
			} else {
				fmt.Println("Task deleted successfully.")
			}

			fmt.Println("Press Enter to continue...")
			reader.ReadString('\n')
		case "weekly":
			fmt.Println("Running weekly review...")
			todo.ReviewWeekly(todos)
			fmt.Print("Press Enter to continue...")
			_, _ = reader.ReadString('\n')

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

func printProjectHelp() {
	fmt.Println(`Available commands (Project REPL):
  todo       - Open the TODO REPL for this project (exits Project REPL)
  add-todo   - Add a new TODO to this project
  edit-todo  - Edit a TODO item in this project
  delete-todo- Delete a TODO item in this project
  weekly     - Run a weekly review of project tasks
  implement  - Implement a todo
  finish     - Mark a todo as complete
  edit       - Edit project info (tags, notes, name, alias, project type)
  exit       - Exit the Project REPL`)
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

	// Attempt to automatically select the todo if exactly one is ongoing.
	selectedIndex := -1
	ongoingCount := 0
	for i, t := range todos {
		if t.Ongoing {
			ongoingCount++
			selectedIndex = i
		}
	}

	// If exactly one ongoing task exists, prompt user whether to proceed.
	if ongoingCount == 1 {
		fmt.Printf("Automatically selecting the only ongoing task: %d. Proceed? (y/n): ", selectedIndex+1)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			// User did not want to use the automatic selection; prompt for manual input.
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
			selectedIndex = index - 1
		}
	} else {
		// If there isn't exactly one ongoing task, prompt the user.
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
		selectedIndex = index - 1
	}

	// Update the selected todo with the new status.
	if err := service.EditTodo(selectedIndex, "", "", status); err != nil {
		fmt.Printf("Error updating todo: %v\n", err)
		return
	}

	fmt.Printf("Todo updated successfully to %s!\n", status)
	description := todos[selectedIndex].Description
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
