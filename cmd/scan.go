package cmd

import (
	"github.com/spf13/cobra"
	"github.com/jacobmaniscalco/goprobe/internal/scan/ble"
)

var macAddress string

// scanCmd represents the scan command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform a network scan to detect vulnerabilities using Nmap.",
	Long: "Scan the specified targets for open ports, services, and " +
	      "potential vulnerabilities using the Nmap library.\nThis command leverages the power " + 
              "of Nmap, a widely used network scanning tool, to perform comprehensive scans on " + 
              "one or more targets.\nThe scan can detect open ports,identify running services, " + 
              "and perform OS detection, among other features.\n Use this command to gather " + 
              "detailed information about the network topology and identify potential security " +
              "weaknesses in your environment.",

	Run: func(cmd *cobra.Command, args []string) {
		ble.StartSniffer(macAddress)
	},
}

func init() {

	ScanCmd.Flags().StringVarP(&macAddress, "macAddress", "m", "", "Specify the MAC address of the device you are attempting to capture.")
}
