package main

import (
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var (
	adb *Adb
)

// Adb is the entrance of adb tool
type Adb struct {
	Path string
}

// NewAdb return a new adb object
func NewAdb() *Adb {
	if adb == nil {
		adb = &Adb{
			Path: "./tools/scrcpy",
		}
	}
	return adb
}

// ListDevices list all android devices
func (adb *Adb) ListDevices() ([]string, error) {
	deviceList := make([]string, 0)

	// Run list devices command
	result, err := adb.runCommand("devices")
	if err != nil {
		return deviceList, fmt.Errorf("list devices failed: %w", err)
	}

	// Parse the devices command's output
	result = strings.ReplaceAll(result, "List of devices attached", "")
	splitList := strings.Split(strings.TrimSpace(result), "\n")

	// Remove the tail
	for _, split := range splitList {
		if split == "" {
			continue
		}

		splitDevice := strings.Split(split, "\t")
		if len(splitDevice) == 0 {
			continue
		}
		deviceList = append(deviceList, strings.TrimSpace(splitDevice[0]))
	}

	return deviceList, nil
}

// Connect connect device with specific ip
func (adb *Adb) Connect(ip string) error {
	output, err := adb.runCommand("connect", ip)
	if err != nil {
		return fmt.Errorf("connect device failed: %w", err)
	}

	if strings.Index(output, "connected to") == -1 {
		return fmt.Errorf("connect device failed, invalid output: %s", output)
	}
	return nil
}

// Disconnect disconnect all devices
func (adb *Adb) Disconnect() error {
	output, err := adb.runCommand("disconnect")
	if err != nil {
		return fmt.Errorf("disconnect failed: %w", err)
	}
	if strings.Index(output, "disconnected everything") == -1 {
		return fmt.Errorf("disconnect failed, invalid output: %s", output)
	}
	return nil
}

// EnableTCPMode enable device's adb tcp mode
func (adb *Adb) EnableTCPMode(devices ...string) error {
	args := make([]string, 0)
	if len(devices) > 0 {
		args = append(args, "-s", devices[0])
	}
	args = append(args, "tcpip", "5555")

	output, err := adb.runCommand(args...)
	if err != nil {
		return fmt.Errorf("enable tcp mod failed: %w", err)
	}
	if strings.Index(output, "restarting in TCP mode port: 5555") == -1 {
		return fmt.Errorf("enable tcp mod failed, invalid output: %s", output)
	}
	return nil
}

// GetWlanIP get the wifi ip
func (adb *Adb) GetWlanIP(devices ...string) (string, error) {
	args := make([]string, 0)
	if len(devices) > 0 {
		args = append(args, "-s", devices[0])
	}
	args = append(args, "shell", "ip -f inet addr show wlan0")
	output, err := adb.runCommand(args...)
	if err != nil {
		return "", fmt.Errorf("get wlan ip failed: %w", err)
	}

	result := regexp.MustCompile("inet (.*)/").FindAllStringSubmatch(output, -1)
	if result == nil || len(result) <= 0 || len(result[0]) < 2 {
		return "", fmt.Errorf("get wlan ip failed, invalid output: %s", output)
	}
	return result[0][1], nil
}

// runCommand run the specific adb command
func (adb *Adb) runCommand(args ...string) (string, error) {
	var command string
	if adb.Path != "" {
		command = path.Join(adb.Path, "adb")
	} else {
		command = "adb"
	}
	output, err := exec.Command(command, args...).Output()
	return string(output), err
}
