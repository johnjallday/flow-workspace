package project

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// editProjectInfo loads the project metadata from the given filename,
// allows the user to interactively edit the name, alias, project type, notes, and tags,
// and then saves the changes back to the file.
func editProjectInfo(filename string) error {
	// Load the project metadata.
	proj, err := LoadProjectInfo(filename)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		// Display current project info.
		fmt.Println("Current Project Info:")
		fmt.Println("1) Name         :", proj.Name)
		fmt.Println("2) Alias        :", proj.Alias)
		fmt.Println("3) Project Type :", proj.ProjectType)
		fmt.Println("4) Notes        :", strings.Join(proj.Notes, ", "))
		fmt.Println("5) Tags         :", strings.Join(proj.Tags, ", "))
		fmt.Println("6) Finish editing")
		fmt.Print("Enter option number to edit: ")

		option, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading option: %v", err)
		}
		option = strings.TrimSpace(option)

		if option == "6" {
			break
		}

		switch option {
		case "1":
			fmt.Print("Enter new name: ")
			newVal, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading name: %v", err)
			}
			proj.Name = strings.TrimSpace(newVal)
		case "2":
			fmt.Print("Enter new alias: ")
			newVal, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading alias: %v", err)
			}
			proj.Alias = strings.TrimSpace(newVal)
		case "3":
			fmt.Print("Enter new project type: ")
			newVal, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading project type: %v", err)
			}
			proj.ProjectType = strings.TrimSpace(newVal)
		case "4":
			fmt.Print("Enter new notes (comma separated): ")
			newVal, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading notes: %v", err)
			}
			proj.Notes = parseList(newVal)
		case "5":
			fmt.Print("Enter new tags (comma separated): ")
			newVal, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading tags: %v", err)
			}
			proj.Tags = parseList(newVal)
		default:
			fmt.Println("Invalid option. Please choose a valid number.")
		}
		fmt.Println()
	}

	// Update the modification timestamp.
	proj.DateModified = time.Now()

	// Save updated project info back to the file.
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(proj); err != nil {
		return fmt.Errorf("failed to encode project info: %v", err)
	}
	fmt.Println("Project info updated successfully.")
	return nil
}

// parseList splits a comma-separated string into a slice of trimmed strings.
func parseList(input string) []string {
	parts := strings.Split(input, ",")
	var list []string
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			list = append(list, item)
		}
	}
	return list
}
