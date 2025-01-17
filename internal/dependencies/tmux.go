package dependencies

import (
	"fmt"
	"os/exec"
)

type TmuxManager struct {
}

func NewTmuxManager() *TmuxManager {
	return &TmuxManager{}
}

// Check if the tmux session exists
func (t *TmuxManager) SessionExists(sessionName string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", sessionName)
	err := cmd.Run()
	return err == nil
}

// Send a command to a tmux session
func (t *TmuxManager) SendCommand(sessionName, command string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", sessionName, command, "C-m")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to send command to tmux session: %w", err)
	}
	return nil
}
