package main

import (
	"time"
	"fmt"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

var test = bluetooth.ParseMAC

func main() {
  	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	address, err := adapter.Address()
	if err != nil {
		fmt.Println("error getting address.")
	}

	fmt.Println(address)


  	// Define the peripheral device info.
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "Go Bluetooth",
  	}))
  
  	// Start advertising
	must("start adv", adv.Start())

	println("advertising...")
	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
