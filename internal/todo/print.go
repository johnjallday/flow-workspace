// internal/todo/print.go
package todo

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

// DisplayTodos prints the list of TODOs in a formatted, colorized table.
// The output adapts based on the terminal width.
func PrintTodos(todos []Todo) {
	// Determine the terminal width.
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80 // fallback width
	}

	// Choose display mode based on terminal width.
	var mode string
	switch {
	case width < 80:
		mode = "compact"
	case width < 120:
		mode = "medium"
	default:
		mode = "full"
	}

	if len(todos) == 0 {
		fmt.Println("No TODOs found.")
		return
	}

	// Create a new tablewriter table.
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false) // No border lines, similar to tabwriter output.

	// Define color styles.
	completeStatus := color.New(color.FgGreen)
	incompleteStatus := color.New(color.FgRed)
	descIncomplete := color.New(color.FgWhite, color.Bold)
	descComplete := color.New(color.FgHiBlack)

	// Set the table header based on the chosen mode.
	switch mode {
	case "compact":
		table.SetHeader([]string{"No.", "Status", "Description"})
	case "medium":
		table.SetHeader([]string{"No.", "Status", "Description", "Due Date"})
	case "full":
		table.SetHeader([]string{"No.", "Status", "Description", "Due Date", "Created", "Project", "Workspace"})
	}

	// Append each TODO as a row in the table.
	for i, t := range todos {
		// Choose a colored status string.
		var status string
		if !t.CompletedDate.IsZero() {
			status = completeStatus.Sprint("[x]")
		} else {
			status = incompleteStatus.Sprint("[ ]")
		}

		// Colorize the description differently if the task is completed.
		var desc string
		if !t.CompletedDate.IsZero() {
			desc = descComplete.Sprint(t.Description)
		} else {
			desc = descIncomplete.Sprint(t.Description)
		}

		// Format dates if set.
		dueDate := ""
		if !t.DueDate.IsZero() {
			dueDate = t.DueDate.Format("2006-01-02")
		}
		created := ""
		if !t.CreatedDate.IsZero() {
			created = t.CreatedDate.Format("2006-01-02")
		}

		// Build the row based on the display mode.
		rowNo := fmt.Sprintf("%d", i+1)
		switch mode {
		case "compact":
			table.Append([]string{rowNo, status, desc})
		case "medium":
			table.Append([]string{rowNo, status, desc, dueDate})
		case "full":
			table.Append([]string{rowNo, status, desc, dueDate, created, t.ProjectName, t.WorkspaceName})
		}
	}

	// Render the table.
	table.Render()
}
