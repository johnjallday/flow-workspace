package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// PrintTree recursively prints an ASCII tree of the given directory.
// It ignores any entry whose name is ".git".
func PrintTree(currentDir, indent string) {
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Printf("%s[ERROR reading directory: %v]\n", indent, err)
		return
	}

	// Filter out entries that should be ignored (e.g. ".git")
	var filteredEntries []os.DirEntry
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		filteredEntries = append(filteredEntries, entry)
	}

	// Iterate over the filtered entries to print each one.
	for i, entry := range filteredEntries {
		isLast := i == len(filteredEntries)-1

		var connector string
		if isLast {
			connector = "└── "
		} else {
			connector = "├── "
		}

		fmt.Printf("%s%s%s\n", indent, connector, entry.Name())

		if entry.IsDir() {
			var newIndent string
			if isLast {
				newIndent = indent + "    "
			} else {
				newIndent = indent + "│   "
			}
			PrintTree(filepath.Join(currentDir, entry.Name()), newIndent)
		}
	}
}
