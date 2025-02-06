// internal/todo/print.go
package todo

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// DisplayTodos prints the list of TODOs in a formatted, colorized table.
// The output adapts based on the terminal width.
func DisplayTodos(todos []Todo) {
	// Attempt to determine the terminal width.
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Fallback to a default width if unable to get size.
		width = 80
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

	// Create a tabwriter for aligned output.
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)

	// Define color styles.
	headerStyle := color.New(color.FgHiGreen, color.Bold)
	completeStatus := color.New(color.FgGreen)
	incompleteStatus := color.New(color.FgRed)
	descIncomplete := color.New(color.FgWhite, color.Bold)
	descComplete := color.New(color.FgHiBlack)

	// Print header and separator based on the chosen mode.
	switch mode {
	case "compact":
		headerStyle.Fprintln(writer, "No.\tStatus\tDescription")
		headerStyle.Fprintln(writer, "----\t------\t-----------")
	case "medium":
		headerStyle.Fprintln(writer, "No.\tStatus\tDescription\tDue Date")
		headerStyle.Fprintln(writer, "----\t------\t-----------\t--------")
	case "full":
		headerStyle.Fprintln(writer, "No.\tStatus\tDescription\tDue Date\tCreated\tProject\tWorkspace")
		headerStyle.Fprintln(writer, "----\t------\t-----------\t--------\t-------\t-------\t---------")
	}

	// Iterate over each TODO item.
	for i, t := range todos {
		// Choose a colored status string.
		var status string
		if t.Completed {
			status = completeStatus.Sprint("[x]")
		} else {
			status = incompleteStatus.Sprint("[ ]")
		}

		// Colorize the description differently if the task is completed.
		var desc string
		if t.Completed {
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

		// Print a row based on the mode.
		switch mode {
		case "compact":
			fmt.Fprintf(writer, "%d\t%s\t%s\n", i+1, status, desc)
		case "medium":
			fmt.Fprintf(writer, "%d\t%s\t%s\t%s\n", i+1, status, desc, dueDate)
		case "full":
			fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
				i+1, status, desc, dueDate, created, t.ProjectName, t.WorkspaceName)
		}
	}

	// Flush the writer to output the formatted table.
	writer.Flush()
}
