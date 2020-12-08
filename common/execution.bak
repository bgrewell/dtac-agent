package common

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
	"time"
)

func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}

func ExecutePipedCmds(commands []string) (output string, err error) {
	if len(commands) < 2 {
		return "", fmt.Errorf("you must pass 2 or more commands to pipe them")
	}
	cmds := make([]*exec.Cmd, len(commands))
	for idx := 0; idx < len(commands); idx++ {
		// break the command into it's fields
		cmdParts := strings.Fields(commands[idx])
		exe, err := exec.LookPath(cmdParts[0])
		if err != nil {
			return "", err
		}
		// setup execution
		cmd := exec.Command(exe, cmdParts[1:]...)
		cmds[idx] = cmd
	}

	bytesout, stderr, err := Pipeline(cmds...)
	if err != nil {
		return "", err
	}
	if len(stderr) > 0 && string(stderr) != "" {
		log.WithFields(log.Fields{
			"stdout": string(bytesout),
			"stderr": string(stderr),
			"err": err,

		}).Debug("command execution had stderr output")
	}
	return string(bytesout), err
}

// ExecuteCmds executes a slice of commands and returns the output and errors from each
func ExecuteCmds(commands []string) (outputs []string, errs []error) {
	outputs = make([]string, len(commands))
	errs = make([]error, len(commands))
	for idx, command := range commands {
		outputs[idx], errs[idx] = ExecuteCmd(command)
	}
	return outputs, errs
}

// ExecuteCmd executes commands and returns the output and any errors
func ExecuteCmd(command string) (output string, err error) {
	cmdParts := strings.Split(command, " ")
	exename, err := exec.LookPath(cmdParts[0])
	exe := exec.Command(exename, cmdParts[1:]...)
	out, err := exe.CombinedOutput()
	cmdLogger := log.WithFields(log.Fields{
		"command": command,
		"path":    exename,
		"args":    cmdParts[1:],
		"output":  string(out),
		"err":     err,
	})
	if err != nil {
		cmdLogger.Error("command executed with errors")
	} else {
		cmdLogger.Debug("command executed successfully")
	}
	return string(out), err
}

// ExecuteCmdWithTimeout executes commands with a timeout. If the timeout occurs the command is terminated and an error is returned
func ExecuteCmdWithTimeout(command string, seconds int) (output string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()
	cmdParts := strings.Split(command, " ")
	exename, err := exec.LookPath(cmdParts[0])
	outBytes, err := exec.CommandContext(ctx, exename, cmdParts[1:]...).Output()
	if ctx.Err() == context.DeadlineExceeded {
		err = errors.New("command execution timeout exceeded")
	}
	return string(outBytes), err
}
