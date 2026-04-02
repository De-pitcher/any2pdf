package utils

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// CommandExists checks if a command is available in PATH
func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// GetCommandPath returns the full path to a command
func GetCommandPath(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("command not found: %s", name)
	}
	return path, nil
}

// ExecOptions contains options for command execution
type ExecOptions struct {
	Timeout time.Duration
	Quiet   bool
}

// ExecResult contains the result of a command execution
type ExecResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
}

// ExecuteCommand runs a command with the given arguments and options
func ExecuteCommand(name string, args []string, opts ExecOptions) ExecResult {
	// Create context with timeout if specified
	var ctx context.Context
	var cancel context.CancelFunc
	
	if opts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	// Create command
	cmd := exec.CommandContext(ctx, name, args...)

	// Capture output
	stdout, err := cmd.Output()
	
	result := ExecResult{
		Stdout: string(stdout),
	}

	if err != nil {
		result.Error = err
		
		// Extract exit code and stderr if available
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			result.Stderr = string(exitErr.Stderr)
		}
	}

	return result
}

// BuildCommand creates a command string for display purposes
func BuildCommand(name string, args []string) string {
	cmd := name
	for _, arg := range args {
		cmd += " " + arg
	}
	return cmd
}
