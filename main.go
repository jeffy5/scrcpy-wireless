package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	adb = NewAdb()
	// Disconnect all devices
	err := adb.Disconnect()
	if err != nil {
		panic(err)
	}

	// List all devices
	deviceList, err := adb.ListDevices()
	if err != nil {
		panic(err)
	}

	// Return if device list is empty
	if len(deviceList) == 0 {
		fmt.Println("Not found devices...")
		return
	}

	// Start scrcpy directly if only one device has connected
	if len(deviceList) == 1 {
		err = startScrcpyWithDevice(deviceList[0])
		if err != nil {
			panic(err)
		}
		return
	}

	// Let user pick which device want to start
	selected := pickDevice(deviceList)
	for selected < 1 || selected > len(deviceList) {
		selected = pickDevice(deviceList, struct{}{})
	}

	// Start scrcpy
	startScrcpyWithDevice(deviceList[selected-1])
}

func pickDevice(deviceList []string, isInvalid ...struct{}) int {
	exec.Command("cmd.exe", "/c", "cls")
	if len(isInvalid) > 0 {
		fmt.Println("Invalid selected device")
	}
	fmt.Println("Please select which device you want to start:")
	for i, device := range deviceList {
		fmt.Printf("%d. %s\n", i, device)
	}

	cmdReader := bufio.NewReader(os.Stdin)
	input, err := cmdReader.ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("read command input failed: %w", err))
	}
	selected, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	return selected
}

func startScrcpyWithDevice(device string) error {
	scrcpy = NewScrcpy()

	// Get device's ip
	ip, err := adb.GetWlanIP(device)
	if err != nil {
		return err
	}

	// Enable tcp mode
	err = adb.EnableTCPMode(device)
	if err != nil {
		return err
	}

	// Connect to device with specific ip
	err = adb.Connect(ip)
	if err != nil {
		return err
	}

	// Start scrcpy
	err = scrcpy.Start(fmt.Sprintf("%s:5555", ip))
	if err != nil {
		return err
	}
	return nil
}
