package todo

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/todo"
)

func Complete(todoFile string, reader *bufio.Reader) {
	tasks, err := todo.LoadAllTodos(todoFile)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks available to complete.")
		return
	}

	fmt.Println("Incomplete Tasks:")
	var incompleteTasks []todo.Todo
	for _, task := range tasks { // renamed loop variable to "task"
		if !task.Completed {
			incompleteTasks = append(incompleteTasks, task)
		}
	}

	if len(incompleteTasks) == 0 {
		fmt.Println("All tasks are already completed.")
		return
	}

	for i, task := range incompleteTasks { // renamed loop variable
		dueDate := "No due date"
		if !task.DueDate.IsZero() {
			dueDate = task.DueDate.Format("2006-01-02")
		}
		project := "No project"
		if task.ProjectName != "" {
			project = task.ProjectName
		}
		workspace := "No workspace"
		if task.WorkspaceName != "" {
			workspace = task.WorkspaceName
		}
		fmt.Printf("%d. %s (Due: %s, Project: %s, Workspace: %s)\n",
			i+1, task.Description, dueDate, project, workspace)
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
	for i, task := range tasks { // renamed loop variable
		if task.Description == selectedTask.Description &&
			task.CreatedDate.Equal(selectedTask.CreatedDate) &&
			task.DueDate.Equal(selectedTask.DueDate) &&
			task.ProjectName == selectedTask.ProjectName &&
			task.WorkspaceName == selectedTask.WorkspaceName &&
			!task.Completed {
			tasks[i].Completed = true
			break
		}
	}

	var updatedContent string
	updatedContent += "# todo\n\n"
	for _, task := range tasks { // renamed loop variable
		status := "[ ]"
		if task.Completed {
			status = "[x]"
		}
		taskLine := fmt.Sprintf("- %s %s", status, task.Description)
		if !task.CreatedDate.IsZero() {
			taskLine += fmt.Sprintf(" #created:%s", task.CreatedDate.Format("2006-01-02"))
		}
		if !task.DueDate.IsZero() {
			taskLine += fmt.Sprintf(" #due:%s", task.DueDate.Format("2006-01-02"))
		}
		if task.ProjectName != "" {
			taskLine += fmt.Sprintf(" #project:%s", task.ProjectName)
		}
		if task.WorkspaceName != "" {
			taskLine += fmt.Sprintf(" #workspace:%s", task.WorkspaceName)
		}
		taskLine += "\n"
		updatedContent += taskLine
	}

	if err := todo.BackupFile(todoFile); err != nil {
		fmt.Printf("Failed to backup 'todo.md': %v\n", err)
	}

	if err := todo.WriteFileContent(todoFile, updatedContent); err != nil {
		fmt.Println("Error writing to 'todo.md':", err)
		return
	}

	fmt.Println("Task marked as completed successfully.")
}
