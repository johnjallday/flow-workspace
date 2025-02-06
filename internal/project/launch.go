package project

import (
	"fmt"
	"os/exec"
	"strings"
)

// LaunchSession creates a new tmux session for the project.
// It sets the session's working directory to projectDir, splits the window horizontally,
// and in the right pane it opens "todo.md" using vim.
func LaunchProject(projectDir string, sessionName string) error {
	// 1. Create a new detached tmux session with working directory set to projectDir.
	newSessionCmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", projectDir)
	if output, err := newSessionCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create tmux session: %v - %s", err, strings.TrimSpace(string(output)))
	}

	// 2. Split the window horizontally with the right pane taking 50 columns.
	// The target is the first window of the session (i.e. sessionName:0).
	splitCmd := exec.Command("tmux", "split-window", "-h", "-l", "50", "-t", sessionName+":0")
	if output, err := splitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to split tmux window: %v - %s", err, strings.TrimSpace(string(output)))
	}

	// 3. In the right pane (pane 1), send the command to open "todo.md" in vim.
	// You can change "vim" to your preferred editor.
	sendKeysCmd := exec.Command("tmux", "send-keys", "-t", sessionName+":0.1", "nvim todo.md", "C-m")
	if output, err := sendKeysCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to send keys to tmux pane: %v - %s", err, strings.TrimSpace(string(output)))
	}

	// 4. Optionally, attach to the new tmux session so the user sees it.
	attachCmd := exec.Command("tmux", "attach-session", "-t", sessionName)
	attachCmd.Stdin = nil
	attachCmd.Stdout = nil
	attachCmd.Stderr = nil

	// Using Run() here will attach and block until the session is detached.
	if err := attachCmd.Run(); err != nil {
		return fmt.Errorf("failed to attach to tmux session: %v", err)
	}

	return nil
}
