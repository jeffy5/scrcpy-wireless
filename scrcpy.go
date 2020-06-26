package main

import (
	"fmt"
	"os/exec"
	"path"
)

var (
	scrcpy *Scrcpy
)

// Scrcpy is the entrance of scrcpy tool
type Scrcpy struct {
	Path string
}

// NewScrcpy return a new scrcpy object
func NewScrcpy() *Scrcpy {
	if scrcpy == nil {
		scrcpy = &Scrcpy{
			Path: "./tools/scrcpy",
		}
	}
	return scrcpy
}

// Start scrcpy
func (scrcpy *Scrcpy) Start(devices ...string) error {
	args := make([]string, 0)
	if len(devices) > 0 {
		args = append(args, "-s", devices[0])
	}

	_, err := scrcpy.runCommand(args...)
	if err != nil {
		return fmt.Errorf("start scrcpy failed: %w", err)
	}
	return nil
}

// runCommand run the specific adb command
func (scrcpy *Scrcpy) runCommand(args ...string) (string, error) {
	var command string
	if scrcpy.Path != "" {
		command = path.Join(scrcpy.Path, "scrcpy")
	} else {
		command = "scrcpy"
	}
	output, err := exec.Command(command, args...).Output()
	return string(output), err
}
