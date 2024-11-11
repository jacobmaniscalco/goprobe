package main

import (
	"fmt"
	"tinygo.org/x/bluetooth"
//	"time"
)

var adapter = bluetooth.DefaultAdapter

func main() {

	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	// Start scanning.
	fmt.Println("scanning...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {

		if device.LocalName() == "Go Bluetooth" {
			adapter.StopScan()
			fmt.Println("found device:", device.Address.String(), device.RSSI, device.LocalName())
			connectToDevice(adapter, device.Address)
		}
	})
	must("start scan", err)
}

func connectToDevice(adapter *bluetooth.Adapter, address bluetooth.Address) {

	device, err := adapter.Connect(address, bluetooth.ConnectionParams{})
	if err != nil {
		fmt.Printf("Error connecting to device: %v\n", err)
	}
	fmt.Printf("Device Address: %v\n", device.Address)

	// ensure connection is stable
	//time.Sleep(5 * time.Second)

	serviceUUID := bluetooth.UUID([4]uint32{1,2,3,4})
	characteristicUUID := bluetooth.UUID([4]uint32{5,6,7,8})

	services, err := device.DiscoverServices([]bluetooth.UUID{serviceUUID})
	if err != nil {
		fmt.Println("Error discovering service:", err)
	}

	if len(services) == 0 {
		panic("could not find any services")
	}

	for _, service := range services {
		characteristics, err := service.DiscoverCharacteristics([]bluetooth.UUID{characteristicUUID})

		if err != nil {
			fmt.Println("Error reading characteristics from service: ", err)
		}

		for _, characteristic := range characteristics {

			data := make([]byte, 128)

			n, err := characteristic.Read(data)

			if err != nil {
				fmt.Println("Error reading characteristic: ", err)
				continue
			}

			fmt.Printf("Read %d bytes from characteristic: %s\n", n, string(data[:n]))
		}

	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
