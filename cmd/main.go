package main

import (
	"github.com/johnjallday/flow-workspace/internal/repl"
)

func main() {
	// Define the root workspace directory
	rootWorkspaceDir := "/Users/jj/Workspace/"

	// Verify that the root directory exists
	// (Already handled in REPL's selectWorkspace function)

	// Start the REPL with the root workspace directory
	repl.StartREPL(rootWorkspaceDir)
}
