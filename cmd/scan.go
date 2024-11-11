package cmd

import (
	//"fmt"

	//"github.com/charmbracelet/bubbletea"
	"github.com/jacobmaniscalco/goprobe/internal/scan"
	//"github.com/jacobmaniscalco/goprobe/internal/scan/snmp"
	//"github.com/jacobmaniscalco/goprobe/ui/components"
	"github.com/spf13/cobra"
	"github.com/jacobmaniscalco/goprobe/internal/scan/ble"
)

var scanOptions scan.ScanOptions
var ports string
var iface string

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

		ble.BLE_scan(iface)
		
		//snmp.ScanSNMPDevice(scanOptions, 161)

		//p := tea.NewProgram(components.NewModel(scanOptions))

		//finalModel, err := p.Run()

		//if err != nil {
		//	fmt.Println("Error with bubbletea rendering: ", err)
		//}

		//if finalModel, ok := finalModel.(components.ScannerModel); ok {

		//	if finalModel.Err != nil {
		//		fmt.Println("Error running scan: ", finalModel.Err)
		//	}

//			if finalModel.Result != "" {
//				fmt.Println(finalModel.Result)
//			}
//		}
	},
}

func init() {

	ScanCmd.Flags().StringVarP(&iface, "interface", "i", "", "Specify interface")

	
	ScanCmd.Flags().StringVarP(&scanOptions.Host, "target", "t", "",
	"Specify the target IP address or range of IP addresses to be scanned.\n" + 
	"This can be a single IP, a subnet, or a list of IPs.\n")

	ScanCmd.Flags().StringVarP(&ports, "ports", "p", "", 
	"Define which ports to scan on the target. You can specify a single port, \n" +
	"a range of ports (e.g., 1-1024), or multiple ports separated by commas.\n")

	ScanCmd.Flags().StringVarP(&scanOptions.Script, "script", "s", "",
	"Indicate the Nmap scripts to be used during the scan.\n" +
	"You can provide a specific script or a category of scripts (e.g., default, vuln).\n")

	ScanCmd.Flags().StringVarP(&scanOptions.Timing, "timing", "", "",
	"Set the timing template for the scan, which controls the speed\n " +
	"and stealthiness of the scan. Timing templates range from 0\n " +
	"(paranoid) to 5 (insane), with higher numbers increasing speed\n " +
	"but also the risk of detection.\n")

	ScanCmd.Flags().BoolVar(&scanOptions.SkipHostDiscovery, "skip-host-discovery",false,
	"Disable host discovery phase, which skips the process of \n" + 
	"identifying live hosts and assumes that the targets are up.\n")

	ScanCmd.Flags().BoolVarP(&scanOptions.Aggressive, "aggressive", "a", false,
	"Enable aggressive scan options, which include service detection, \n" +
	"OS detection, and additional scanning techniques to gather more information.\n")

	ScanCmd.Flags().BoolVar(&scanOptions.ServiceDetection, "service-detection", false,
	"Enable detection of services running on the open \n" + 
	"ports of the target.This helps in identifying the applications " +
	"and versions running on the target.\n")

	ScanCmd.Flags().BoolVar(&scanOptions.OsDetection, "os-detection", false,
	"Enable OS detection to determine the operating system \n" + 
	"of the target host. This uses various techniques to guess the \n" + 
	"OS based on network responses.\n")

	ScanCmd.Flags().StringVar(&scanOptions.ScanType, "scan-type", "", 
	"Specify scan type (e.g., sS for SYN, sT for TCP connect, sU for UDP)\n")
}
