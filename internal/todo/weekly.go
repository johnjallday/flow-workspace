package todo

import (
	"fmt"
	"time"
)

// ReviewWeekly filters todos completed this week (Saturday to Friday)
// and prints a report of both completed and unfinished tasks.
// It also indicates if today is the review day (Friday).
func ReviewWeekly(todos []Todo) {
	now := time.Now()
	// Reset today's time to midnight to simplify date comparisons.
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Compute a custom weekday index where Saturday is 0, Sunday is 1, ..., Friday is 6.
	customWeekday := (int(now.Weekday()) + 1) % 7
	// The review period starts on the most recent Saturday.
	reviewStart := today.AddDate(0, 0, -customWeekday)
	// The review period ends on Friday (6 days after Saturday).
	reviewEnd := reviewStart.AddDate(0, 0, 6)

	fmt.Printf("Review Period: %s to %s\n", reviewStart.Format("2006-01-02"), reviewEnd.Format("2006-01-02"))

	var completed []Todo
	var unfinished []Todo

	// Iterate over todos and classify them.
	for _, todo := range todos {
		// A non-zero CompletedDate means the task is complete.
		if !todo.CompletedDate.IsZero() {
			// Check if the task was completed within the review period.
			if !todo.CompletedDate.Before(reviewStart) && !todo.CompletedDate.After(reviewEnd) {
				completed = append(completed, todo)
			}
		} else {
			// The task is not yet completed.
			unfinished = append(unfinished, todo)
		}
	}

	fmt.Println("\nCompleted Todos for this week:")
	for _, t := range completed {
		fmt.Printf(" - %s (Completed on: %s, Project: %s, Workspace: %s)\n",
			t.Description,
			t.CompletedDate.Format("2006-01-02"),
			t.ProjectName,
			t.WorkspaceName)
	}

	fmt.Println("\nUnfinished Todos:")
	for _, t := range unfinished {
		fmt.Printf(" - %s (Due: %s, Project: %s, Workspace: %s)\n",
			t.Description,
			t.DueDate.Format("2006-01-02"),
			t.ProjectName,
			t.WorkspaceName)
	}

	// Check if today is the review day (Friday).
	if now.Weekday() == time.Friday {
		fmt.Println("\nToday is review day!")
	} else {
		fmt.Println("\nToday is not review day.")
	}
}
