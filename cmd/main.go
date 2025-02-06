package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	projectTodo "github.com/johnjallday/flow-workspace/internal/project/todo"
	"github.com/johnjallday/flow-workspace/internal/repl"
	"github.com/johnjallday/flow-workspace/internal/root"
	workspaceTodo "github.com/johnjallday/flow-workspace/internal/workspace/todo"
)

func main() {
	flag.Parse()
	args := flag.Args()

	// If a command is provided, handle it.
	if len(args) > 0 {
		switch args[0] {
		case "todo":
			// Use provided directory or default to current working directory.
			var dir string
			if len(args) >= 2 {
				dir = args[1]
			} else {
				var err error
				dir, err = os.Getwd()
				if err != nil {
					log.Fatalf("Failed to get current directory: %v", err)
				}
			}

			// Check if the directory exists.
			info, err := os.Stat(dir)
			if err != nil {
				log.Fatalf("Directory '%s' does not exist: %v", dir, err)
			}
			if !info.IsDir() {
				log.Fatalf("Provided path '%s' is not a directory.", dir)
			}

			// Clean up the path.
			dir = filepath.Clean(dir)

			// Determine scope by checking for known marker files/folders.
			// 1. Root scope: check for a ".config" folder.
			configPath := filepath.Join(dir, ".config")
			if info, err := os.Stat(configPath); err == nil && info.IsDir() {
				// Run the root-level REPL.
				root.ListAllTodos(dir)
				return
			}

			// 2. Workspace scope: check for a "projects.toml" file.
			projectsTomlPath := filepath.Join(dir, "projects.toml")
			if _, err := os.Stat(projectsTomlPath); err == nil {
				// Run the workspace-level REPL.
				workspaceTodo.StartTodoREPL(dir)
				return
			}

			// 3. Project scope: check for a "project_info.toml" file.
			projectInfoPath := filepath.Join(dir, "project_info.toml")
			if _, err := os.Stat(projectInfoPath); err == nil {
				// Run the project-level (todo) REPL.

				todoFilePath := filepath.Join(dir, "todo.md")
				projectTodo.StartTodoREPL(todoFilePath)
				return
			}

			// Fallback: no known scope marker found. Default to the project-level REPL.
			fmt.Println("No known scope markers found in directory. Running project-level REPL by default.")
			projectTodo.StartTodoREPL(dir)

		default:
			fmt.Println("Unknown command. Available commands: todo")
		}
		return
	}

	// No command provided: start the general scope-detecting REPL.
	repl.StartREPL()
}
