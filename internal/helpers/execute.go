package helpers

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// RunAsUser executes a system command as a specific user and returns stdout and any error.
func RunAsUser(cmd string, user string) (string, error) {
	// Split command
	cmds := strings.Fields(cmd)
	// Construct the command to run with 'sudo' and '-u' to execute as a specific user
	commandArgs := append([]string{"sudo", "-u", user}, cmds...)
	command := exec.Command(commandArgs[0], commandArgs[1:]...)

	// Capture stdout and stderr
	var out bytes.Buffer
	var errBuf bytes.Buffer
	command.Stdout = &out
	command.Stderr = &errBuf

	// Run the command
	err := command.Run()
	if err != nil {
		// Handle exit errors specifically
		if exitError, ok := err.(*exec.ExitError); ok {
			// Return a formatted error with the exit code
			return "", errors.New("process exited with error code " + strconv.Itoa(exitError.ExitCode()) + ": " + errBuf.String())
		}
		// For other errors
		return "", err
	}

	// Return the stdout as a string
	return strings.Trim(out.String(), "\n"), nil
}
