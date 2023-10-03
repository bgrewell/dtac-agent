package common

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/configuration"
)

func PrepForInstall() error {
	err := createDirectories()
	if err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	err = copyBinary()
	if err != nil {
		return fmt.Errorf("failed to copy binary: %v", err)
	}

	fmt.Println("Running install...")
	cmd := exec.Command(configuration.BINARY_NAME, "--install")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command(configuration.BINARY_NAME, "--start")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	return err
}

func copyBinary() error {
	fmt.Println("Preparing Binary...")
	src, err := os.Executable()
	if err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(configuration.BINARY_NAME)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return os.Chmod(configuration.BINARY_NAME, 0755)
}

func createDirectories() error {
	fmt.Println("Creating directories...")
	directories := []string{
		configuration.GLOBAL_CONFIG_LOCATION,
		configuration.DEFAULT_BINARY_LOCATION,
		configuration.DEFAULT_PLUGIN_LOCATION,
		configuration.GLOBAL_CERT_LOCATION,
	}
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func BinaryIsCorrect() bool {
	current, err := os.Executable()
	if err != nil {
		fmt.Printf("failed to get current executable: %v\n", err)
		return false
	}

	return current == configuration.BINARY_NAME
}
