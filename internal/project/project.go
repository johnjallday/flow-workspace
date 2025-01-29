// internal/project/project.go

package project

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	openai "github.com/sashabaranov/go-openai"
)

// Project represents the structure of a single project_info.toml file.
type Project struct {
	Name         string    `toml:"name"`              // The name of the project.
	Alias        string    `toml:"alias"`             // An alternative name or shorthand for the project.
	ProjectType  string    `toml:"project_type"`      // The type/category of the project.
	Tags         []string  `toml:"tags"`              // A list of tags associated with the project.
	DateCreated  time.Time `toml:"date_created"`      // The creation date of the project.
	DateModified time.Time `toml:"date_modified"`     // The last modified date of the project.
	Notes        []string  `toml:"notes"`             // Any additional notes or comments about the project.
	Path         string    `toml:"path"`              // The file system path to the project's root directory.
	GitURL       string    `toml:"git_url,omitempty"` // Optional Git repository URL.
}

// Projects represents a collection of Project entries aggregated into projects.toml.
type Projects struct {
	Projects []Project `toml:"projects"` // A slice of Project structs.
}

// LoadProjectInfo reads and parses a project_info.toml file into a Project struct.
func LoadProjectInfo(filename string) (*Project, error) {
	var proj Project

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file '%s' does not exist", filename)
	}

	// Decode the TOML file into the Project struct
	if _, err := toml.DecodeFile(filename, &proj); err != nil {
		return nil, fmt.Errorf("error decoding TOML file: %w", err)
	}

	// Validate mandatory fields
	if proj.Name == "" {
		return nil, fmt.Errorf("'name' field cannot be empty in '%s'", filename)
	}

	// Set default values for optional fields if empty
	if proj.Alias == "" {
		proj.Alias = proj.Name
	}

	if proj.ProjectType == "" {
		proj.ProjectType = "General" // Default project type
	}

	if proj.Path == "" {
		proj.Path = "./" // Default path
	}

	return &proj, nil
}

// LoadProjectsInfo reads and parses a projects.toml file into a Projects struct.
func LoadProjectsInfo(filename string) (*Projects, error) {
	var projs Projects

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file '%s' does not exist", filename)
	}

	// Decode the TOML file into the Projects struct
	if _, err := toml.DecodeFile(filename, &projs); err != nil {
		return nil, fmt.Errorf("error decoding TOML file: %w", err)
	}

	return &projs, nil
}

// SaveProjectsInfo marshals the Projects struct and writes it to projects.toml.
func SaveProjectsInfo(projs *Projects, filename string) error {
	// Marshal the Projects struct into TOML format
	output, err := toml.Marshal(projs)
	if err != nil {
		return fmt.Errorf("failed to marshal projects to TOML: %w", err)
	}

	// Write the marshaled data to projects.toml
	if err := os.WriteFile(filename, output, 0644); err != nil {
		return fmt.Errorf("failed to write to '%s': %w", filename, err)
	}

	return nil
}

// ListProjects prints all projects in the Projects struct in a formatted manner.
func ListProjects(projs *Projects) {
	if projs == nil || len(projs.Projects) == 0 {
		fmt.Println("No projects found.")
		return
	}

	fmt.Println("\nProjects in this Workspace:")
	for i, proj := range projs.Projects {
		fmt.Printf("--------------------------------------------------\n")
		fmt.Printf("Project #%d\n", i+1)
		fmt.Printf("Name         : %s\n", proj.Name)
		fmt.Printf("Alias        : %s\n", proj.Alias)
		fmt.Printf("Type         : %s\n", proj.ProjectType)
		fmt.Printf("Tags         : %s\n", strings.Join(proj.Tags, ", "))
		fmt.Printf("Date Created : %s\n", proj.DateCreated.Format("January 2, 2006"))
		fmt.Printf("Date Modified: %s\n", proj.DateModified.Format("January 2, 2006"))
		if proj.GitURL != "" {
			fmt.Printf("Git URL      : %s\n", proj.GitURL)
		}
		fmt.Printf("Path         : %s\n", proj.Path)
		fmt.Printf("--------------------------------------------------\n\n")
	}
}

// ScanAndAggregateProjects scans the root directory for project_info.toml files and aggregates them into projects.toml.
func ScanAndAggregateProjects(rootDir string) (*Projects, error) {
	var projects []Project

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Failed to access path %s: %v\n", path, err)
			return nil // Continue walking
		}

		// Check if the file is named 'project_info.toml'
		if !info.IsDir() && info.Name() == "project_info.toml" {
			fmt.Printf("Found project_info.toml: %s\n", path)
			p, err := LoadProjectInfo(path)
			if err != nil {
				fmt.Printf("Failed to load project info from '%s': %v\n", path, err)
				return nil // Continue walking
			}

			projects = append(projects, *p)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the path '%s': %w", rootDir, err)
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("no project_info.toml files found")
	}

	aggregatedProjects := Projects{
		Projects: projects,
	}

	return &aggregatedProjects, nil
}

// OpenInNeovim launches Neovim to edit the specified file.
// If readOnly is true, it opens the file in read-only mode.
func OpenInNeovim(filename string, readOnly bool) error {
	args := []string{filename}
	if readOnly {
		args = append([]string{"-R"}, args...)
	}
	cmd := exec.Command("nvim", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RenderOpenAI reads a template file, calls OpenAI, and returns the GPT output.
func RenderOpenAI(inputContent, templatePath string) (string, error) {
	// 1) Read the template file
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading template file (%s): %w", templatePath, err)
	}

	// 2) Build system/user prompts
	systemPrompt := "Convert the user content into the following template:\n" + string(tmplContent)
	userPrompt := inputContent

	// 3) API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set in environment")
	}

	// 4) Create OpenAI client
	client := openai.NewClient(apiKey)

	// 5) Prepare the chat request
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini20240718, // or GPT-4 if you have access
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
		Temperature: 0,
	}

	// 6) Call OpenAI
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	finalOutput := resp.Choices[0].Message.Content
	return finalOutput, nil
}
