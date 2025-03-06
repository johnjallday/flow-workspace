package todo

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Complete(todoFile string, reader *bufio.Reader) {
	// Load all tasks from the file.
	tasks, err := LoadAllTodos(todoFile)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks available to complete.")
		return
	}

	// Filter out incomplete tasks (i.e. tasks with a zero CompletedDate).
	fmt.Println("Incomplete Tasks:")
	var incompleteTasks []Todo
	for _, task := range tasks {
		if task.CompletedDate.IsZero() {
			incompleteTasks = append(incompleteTasks, task)
		}
	}

	if len(incompleteTasks) == 0 {
		fmt.Println("All tasks are already completed.")
		return
	}

	// Display the incomplete tasks.
	for i, task := range incompleteTasks {
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

	// Prompt the user to select a task to complete.
	fmt.Print("Enter the number of the task to mark as completed: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(incompleteTasks) {
		fmt.Println("Invalid task number.")
		return
	}

	// Find the selected task in the full list and mark it complete by setting CompletedDate.
	selectedTask := incompleteTasks[index-1]
	for i, task := range tasks {
		if task.Description == selectedTask.Description &&
			task.CreatedDate.Equal(selectedTask.CreatedDate) &&
			task.DueDate.Equal(selectedTask.DueDate) &&
			task.ProjectName == selectedTask.ProjectName &&
			task.WorkspaceName == selectedTask.WorkspaceName &&
			task.CompletedDate.IsZero() {
			// Mark the task as completed by setting the CompletedDate to now.
			tasks[i].CompletedDate = time.Now()
			break
		}
	}

	// Rebuild the file content.
	var updatedContent string
	updatedContent += "# todo\n\n"
	for _, task := range tasks {
		status := "[ ]"
		if !task.CompletedDate.IsZero() {
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
		// Add a completed tag if the task is completed.
		if !task.CompletedDate.IsZero() {
			taskLine += fmt.Sprintf(" #completed:%s", task.CompletedDate.Format("2006-01-02"))
		}
		taskLine += "\n"
		updatedContent += taskLine
	}

	// Write the updated content back to the file.
	if err := WriteFileContent(todoFile, updatedContent); err != nil {
		fmt.Println("Error writing to 'todo.md':", err)
		return
	}

	fmt.Println("Task marked as completed successfully.")
}
