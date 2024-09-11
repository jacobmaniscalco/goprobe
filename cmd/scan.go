/*

Copyright © 2024 Jacob Maniscalco
*/
package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
	"github.com/jacobmaniscalco/blue-caterpillar-cli/internal/scan"

)

var target string 
var ports string
var verbose bool

// scanCmd represents the scan command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform a network scan to detect vulnerabilities using Nmap.",
	Long: 
`Scan the specified targets for open ports, services, and potential vulnerabilities using the Nmap library.
This command leverages the power of Nmap, a widely used network scanning tool, to perform comprehensive scans on one or more targets.
The scan can detect open ports, identify running services, and perform OS detection, among other features. 
Use this command to gather detailed information about the network topology and identify potential security weaknesses in your environment.
Examples:
  	- Scan a single IP address:
      	scan 192.168.1.1
  	- Scan a range of IP addresses:
      	scan 192.168.1.1-10
  	- Scan a specific domain:
      	scan example.com
  	- Scan using OS detection:
      	scan --os-detection example.com`,

	Run: func(cmd *cobra.Command, args []string) {	

		_, err := scan.RunScan(target, ports)

		if err != nil {
			fmt.Printf("Error running scan: %v\n", err)
			return
		}
	},
}

func init() {

	// flag for IP address to scan
	ScanCmd.Flags().StringVarP(&target, "target", "t", "", "Target to scan")
	ScanCmd.MarkFlagRequired("target")
        
	// flag for ports to scan
	ScanCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan")

	// flag for verbose output
	ScanCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
