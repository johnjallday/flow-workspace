package todo

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func Complete(todoFile string, reader *bufio.Reader) {
	todos, err := LoadAllTodos(todoFile)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No tasks available to complete.")
		return
	}

	fmt.Println("Incomplete Tasks:")
	var incompleteTasks []Todo
	for _, todo := range todos {
		if !todo.Completed {
			incompleteTasks = append(incompleteTasks, todo)
		}
	}

	if len(incompleteTasks) == 0 {
		fmt.Println("All tasks are already completed.")
		return
	}

	for i, todo := range incompleteTasks {
		dueDate := "No due date"
		if !todo.DueDate.IsZero() {
			dueDate = todo.DueDate.Format("2006-01-02")
		}
		project := "No project"
		if todo.ProjectName != "" {
			project = todo.ProjectName
		}
		workspace := "No workspace"
		if todo.WorkspaceName != "" {
			workspace = todo.WorkspaceName
		}
		fmt.Printf("%d. %s (Due: %s, Project: %s, Workspace: %s)\n",
			i+1, todo.Description, dueDate, project, workspace)
	}

	fmt.Print("Enter the number of the task to mark as completed: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(incompleteTasks) {
		fmt.Println("Invalid task number.")
		return
	}

	selectedTask := incompleteTasks[index-1]
	for i, todo := range todos {
		if todo.Description == selectedTask.Description &&
			todo.CreatedDate.Equal(selectedTask.CreatedDate) &&
			todo.DueDate.Equal(selectedTask.DueDate) &&
			todo.ProjectName == selectedTask.ProjectName &&
			todo.WorkspaceName == selectedTask.WorkspaceName &&
			!todo.Completed {
			todos[i].Completed = true
			break
		}
	}

	var updatedContent string
	updatedContent += "# todo\n\n"
	for _, todo := range todos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		taskLine := fmt.Sprintf("- %s %s", status, todo.Description)
		if !todo.CreatedDate.IsZero() {
			taskLine += fmt.Sprintf(" #created:%s", todo.CreatedDate.Format("2006-01-02"))
		}
		if !todo.DueDate.IsZero() {
			taskLine += fmt.Sprintf(" #due:%s", todo.DueDate.Format("2006-01-02"))
		}
		if todo.ProjectName != "" {
			taskLine += fmt.Sprintf(" #project:%s", todo.ProjectName)
		}
		if todo.WorkspaceName != "" {
			taskLine += fmt.Sprintf(" #workspace:%s", todo.WorkspaceName)
		}
		taskLine += "\n"
		updatedContent += taskLine
	}

	if err := BackupFile(todoFile); err != nil {
		fmt.Printf("Failed to backup 'todo.md': %v\n", err)
	}

	if err := WriteFileContent(todoFile, updatedContent); err != nil {
		fmt.Println("Error writing to 'todo.md':", err)
		return
	}

	fmt.Println("Task marked as completed successfully.")
}
