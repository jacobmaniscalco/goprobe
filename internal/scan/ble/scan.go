package ble


import (
	"bufio"
	"fmt"
	"strings"
	"log"
	"os/exec"
)

type bleScanResults struct {
	source string
	destination string
	rssi string
}

func BLE_scan(iface string) {
	cmd := exec.Command("tshark", "-i", iface, "-T", "fields",
	"-e", "_ws.col.Source",
	"-e", "_ws.col.Destination",
	"-e", "nordic_ble.rssi",
	"-E", "separator=\t",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start tshark: %v", err)
	}
	defer cmd.Process.Kill()


	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		//line = strings.TrimSpace(line)
		fields := strings.Split(line, "\t")

		source := fields[0]
		destination := fields[1]
		rssi := fields[2]

		fmt.Printf("Source: %s\nDestination: %v\nRssi: %s\n\n\n\n", 
			source, destination, rssi)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("tshark command failed: %v", err)
	}
}

func getFirstOrDefault(arr []string) string {
	if len(arr) > 0 {
		return arr[0]
	}

	return "N/A"
}
