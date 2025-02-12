package todo

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/johnjallday/flow-workspace/internal/todo"
)

func Delete(todoFile string, reader *bufio.Reader) {
	todos, err := todo.LoadAllTodos(todoFile)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No tasks available to delete.")
		return
	}

	fmt.Println("All Tasks:")
	for i, todo := range todos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		fmt.Printf("%d. %s %s\n", i+1, status, todo.Description)
	}

	fmt.Print("Enter the number of the task to delete: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(todos) {
		fmt.Println("Invalid task number.")
		return
	}

	todos = append(todos[:index-1], todos[index:]...)

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

	if err := todo.BackupFile(todoFile); err != nil {
		fmt.Printf("Failed to backup 'todo.md': %v\n", err)
	}

	if err := todo.WriteFileContent(todoFile, updatedContent); err != nil {
		fmt.Println("Error writing to 'todo.md':", err)
		return
	}

	fmt.Println("Task deleted successfully.")
}
