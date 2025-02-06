package todo

import (
	"fmt"
	"strings"
)

func CleanUpTodoFile(filename string) error {
	content, err := ReadFileContent(filename)
	if err != nil {
		return fmt.Errorf("error reading '%s': %w", filename, err)
	}

	lines := strings.Split(content, "\n")
	var tasks []string
	var headers []int

	for idx, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# todo") {
			headers = append(headers, idx)
			continue
		}
		if strings.HasPrefix(trimmed, "- [ ]") || strings.HasPrefix(trimmed, "- [x]") {
			tasks = append(tasks, line)
		}
	}

	if len(headers) == 0 {
		lines = append([]string{"# todo", ""}, lines...)
	} else {
		if len(headers) > 1 {
			for i := 1; i < len(headers); i++ {
				headerIdx := headers[i] - i
				if headerIdx >= 0 && headerIdx < len(lines) {
					lines = append(lines[:headerIdx], lines[headerIdx+1:]...)
				}
			}
		}

		linesAfterHeader := lines[0 : headers[0]+1]
		var newTasks []string
		for _, line := range linesAfterHeader {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "- [ ]") || strings.HasPrefix(trimmed, "- [x]") {
				newTasks = append(newTasks, line)
			}
		}
		lines = linesAfterHeader
		tasks = append(newTasks, tasks...)
	}

	updatedLines := []string{"# todo", ""}
	updatedLines = append(updatedLines, tasks...)
	updatedContent := strings.Join(updatedLines, "\n")

	if err := WriteFileContent(filename, updatedContent); err != nil {
		return fmt.Errorf("error writing to '%s': %w", filename, err)
	}

	return nil
}
