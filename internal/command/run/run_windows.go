//go:build windows
// +build windows

package run

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func configureProcess(_ *exec.Cmd) {}

func interruptProcess(cmd *exec.Cmd) error {
	return taskkill(cmd, false)
}

func forceKillProcess(cmd *exec.Cmd) error {
	return taskkill(cmd, true)
}

func taskkill(cmd *exec.Cmd, force bool) error {
	args := []string{"/T", "/PID", strconv.Itoa(cmd.Process.Pid)}
	if force {
		args = append([]string{"/F"}, args...)
	}
	out, err := exec.Command("taskkill", args...).CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(out))
		if message == "" {
			return fmt.Errorf("taskkill: %w", err)
		}
		return fmt.Errorf("taskkill: %w: %s", err, message)
	}
	return nil
}
