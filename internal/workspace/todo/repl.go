package todo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// ----- Minimal definitions duplicated from workspace package ----- //

// Project is a minimal duplicate of the project type used in projects.toml.
type Project struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
	// Add other fields if needed.
}

// Projects represents a collection of Project entries.
type Projects struct {
	Projects []Project `toml:"projects"`
}

// loadProjectsTomlIndependent is a duplicate of workspace.LoadProjectsToml,
// defined here to avoid an import cycle. It loads and decodes the projects.toml
// file from a given workspace directory.
func loadProjectsTomlIndependent(workspacePath string) (*Projects, error) {
	projectsTomlPath := filepath.Join(workspacePath, "projects.toml")
	if _, err := os.Stat(projectsTomlPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file '%s' does not exist", projectsTomlPath)
	}

	var projs Projects
	if _, err := toml.DecodeFile(projectsTomlPath, &projs); err != nil {
		return nil, fmt.Errorf("error decoding TOML file: %w", err)
	}

	return &projs, nil
}

// StartTodoREPL aggregates todos from all projects in the workspace (or root).
// It creates (or reuses) a file named todo_temp.md in the given directory to store
// the aggregated todos. If the file already exists, it opens it without reaggregating.
// In the REPL loop, the user can type "update" to reaggregate and update the file,
// "add" to add a new task (with workspace tag automatically added), or "exit" to quit.
func StartTodoREPL(workspaceDir string) {
	tempFilePath := filepath.Join(workspaceDir, "todo_temp.md")
	var todos []Todo
	var err error

	// Check if todo_temp.md already exists.
	if _, err = os.Stat(tempFilePath); err == nil {
		// File exists: load todos from it.
		todos, err = LoadAllTodos(tempFilePath)
		if err != nil {
			fmt.Printf("Error loading todos from %s: %v\n", tempFilePath, err)
			return
		}
	} else {
		// Otherwise, aggregate todos from all projects in the workspace.
		todos, err = aggregateWorkspaceTodos(workspaceDir)
		if err != nil {
			fmt.Printf("Error aggregating todos: %v\n", err)
			return
		}
		// Write the aggregated todos to todo_temp.md.
		if err := writeAggregatedTodos(tempFilePath, todos); err != nil {
			fmt.Printf("Error writing aggregated todos to %s: %v\n", tempFilePath, err)
			return
		}
	}

	// Start the REPL loop.
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nAggregated TODOs:")
		DisplayTodos(todos)
		fmt.Println("\nAvailable commands: update, exit")
		fmt.Print("Enter command: ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		switch line {
		case "exit":
			return
		case "update":
			// Reaggregate todos and update the temporary file.
			todos, err = aggregateWorkspaceTodos(workspaceDir)
			if err != nil {
				fmt.Printf("Error reaggregating todos: %v\n", err)
				continue
			}
			if err := writeAggregatedTodos(tempFilePath, todos); err != nil {
				fmt.Printf("Error updating file %s: %v\n", tempFilePath, err)
				continue
			}
			fmt.Println("Todos updated.")
		case "add":
			// Prompt the user for a new task.
			fmt.Print("Enter new task description: ")
			desc, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading description: %v\n", err)
				continue
			}
			desc = strings.TrimSpace(desc)
			if desc == "" {
				fmt.Println("Task description cannot be empty.")
				continue
			}

			fmt.Print("Enter due date (YYYY-MM-DD) or leave empty: ")
			dueStr, _ := reader.ReadString('\n')
			dueStr = strings.TrimSpace(dueStr)
			var dueDate time.Time
			if dueStr != "" {
				dt, err := time.Parse("2006-01-02", dueStr)
				if err != nil {
					fmt.Printf("Invalid date format: %v\n", err)
					continue
				}
				dueDate = dt
			}

			// Create a new Todo item.
			newTodo := Todo{
				Description: desc,
				Completed:   false,
				CreatedDate: time.Now(),
				DueDate:     dueDate,
			}
			// Automatically add workspace tag using tagWorkspace.
			newTodo.WorkspaceName = tagWorkspace(workspaceDir)

			// Append the new todo to the list.
			todos = append(todos, newTodo)
			// Update the temporary file.
			if err := writeAggregatedTodos(tempFilePath, todos); err != nil {
				fmt.Printf("Error updating file %s: %v\n", tempFilePath, err)
				continue
			}
			fmt.Println("Task added successfully.")
		default:
			fmt.Println("Unknown command. Available commands: update, add, exit")
		}
	}
}

// aggregateWorkspaceTodos aggregates TODOs for all projects in a workspace.
func aggregateWorkspaceTodos(workspaceDir string) ([]Todo, error) {
	projs, err := loadProjectsTomlIndependent(workspaceDir)
	if err != nil {
		return nil, err
	}

	var aggregated []Todo
	for _, proj := range projs.Projects {
		var projDir string
		if proj.Path != "" && proj.Path != "./" {
			if filepath.IsAbs(proj.Path) {
				projDir = filepath.Clean(proj.Path)
			} else {
				projDir = filepath.Join(workspaceDir, proj.Path)
			}
		} else {
			projDir = filepath.Join(workspaceDir, proj.Name)
		}
		todoFile := filepath.Join(projDir, "todo.md")
		if _, err := os.Stat(todoFile); os.IsNotExist(err) {
			continue
		}
		tasks, err := LoadAllTodos(todoFile)
		if err != nil {
			fmt.Printf("Error loading todos from %s: %v\n", todoFile, err)
			continue
		}
		// Optionally annotate each task with the project name.
		for i := range tasks {
			if tasks[i].ProjectName == "" {
				tasks[i].ProjectName = proj.Name
			}
		}
		aggregated = append(aggregated, tasks...)
	}
	return aggregated, nil
}

// writeAggregatedTodos writes the aggregated todos into the given file.
// It writes a header "# todo" followed by each task line.
func writeAggregatedTodos(filePath string, todos []Todo) error {
	var lines []string
	lines = append(lines, "# todo", "")
	for _, t := range todos {
		status := "[ ]"
		if t.Completed {
			status = "[x]"
		}
		// Build a simple task line.
		line := fmt.Sprintf("- %s %s", status, t.Description)
		if !t.CreatedDate.IsZero() {
			line += fmt.Sprintf(" #created:%s", t.CreatedDate.Format("2006-01-02"))
		}
		if !t.DueDate.IsZero() {
			line += fmt.Sprintf(" #due:%s", t.DueDate.Format("2006-01-02"))
		}
		if t.ProjectName != "" {
			line += fmt.Sprintf(" #project:%s", t.ProjectName)
		}
		if t.WorkspaceName != "" {
			line += fmt.Sprintf(" #workspace:%s", t.WorkspaceName)
		}
		lines = append(lines, line)
	}
	content := strings.Join(lines, "\n")
	return WriteFileContent(filePath, content)
}
