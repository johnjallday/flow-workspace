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
	"github.com/johnjallday/flow-workspace/internal/startup"
	"github.com/johnjallday/flow-workspace/internal/workspace"
)

func main() {
	dbPath := startup.StartDB() // call startDB() from the startup package
	fmt.Println("Database path:", dbPath)

	flag.Parse()
	args := flag.Args()

	// If a command is provided, handle it.
	if len(args) > 0 {
		switch args[0] {
		case "todo":
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
			configPath := filepath.Join(dir, ".config")
			if info, err := os.Stat(configPath); err == nil && info.IsDir() {
				root.ListAllTodos(dir)
				return
			}

			projectsTomlPath := filepath.Join(dir, "projects.toml")
			if _, err := os.Stat(projectsTomlPath); err == nil {
				// Updated to call ListAllTodos instead of StartTodoREPL
				workspace.ListAllTodos(dir)
				return
			}

			projectInfoPath := filepath.Join(dir, "project_info.toml")
			if _, err := os.Stat(projectInfoPath); err == nil {
				todoFilePath := filepath.Join(dir, "todo.md")
				projectTodo.StartTodoREPL(dbPath, todoFilePath)
				return
			}

			// Fallback: no known scope marker found.
			fmt.Println("No known scope markers found in directory.")

		default:
			fmt.Println("Unknown command. Available commands: todo")
		}
		return
	}

	// No command provided: start the general scope-detecting REPL.
	repl.StartREPL(dbPath)
}
