package ble

import (
	"tinygo.org/x/bluetooth"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"fmt"
)

var adapter = bluetooth.DefaultAdapter

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ":" + err.Error())
	}
}

func BLE_scan() {
	//must("enable BLE stack", adapter.Enable())

	//println("scanning...")

	//addr, _ := adapter.Address()

	//fmt.Println(addr)

	//err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
	//	println("found device:", device.Address.String(), device.RSSI, device.LocalName())
	//})
	//must("start scan", err)
	
	interfaceName := "hci0"

	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)

	if err != nil {
		panic(err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println("Captured packet:", packet)
	}

}
