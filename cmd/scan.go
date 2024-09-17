package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/jacobmaniscalco/blue-caterpillar-cli/internal/scan"
	"github.com/jacobmaniscalco/blue-caterpillar-cli/ui/components"
	"github.com/spf13/cobra"
)

var options scan.ScanOptions

// scanCmd represents the scan command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform a network scan to detect vulnerabilities using Nmap.",
	Long: `Scan the specified targets for open ports, services, and 
potential vulnerabilities using the Nmap library.This command leverages the power
of Nmap, a widely used network scanning tool, to perform comprehensive scans on
one or more targets.The scan can detect open ports,identify running services, 
and perform OS detection, among other features. Use this command to gather 
detailed information about the network topology and identify potential security
weaknesses in your environment.`,

	Run: func(cmd *cobra.Command, args []string) {

		p := tea.NewProgram(components.NewModel(options))

		finalModel, err := p.Run()

		if err != nil {
			fmt.Println("Error with bubbletea rendering: ", err)
		}

		if finalModel, ok := finalModel.(components.ScannerModel); ok {

			if finalModel.Err != nil {
				fmt.Println("Error running scan: ", finalModel.Err)
			}

			if finalModel.Result != "" {
				fmt.Println(finalModel.Result)
			}
		}
	},
}

func init() {

	
	ScanCmd.Flags().StringVarP(&options.Target, "target", "t", "",
`Specify the target IP address or range of IP addresses to be scanned. 
This can be a single IP, a subnet, or a list of IPs.
`)
	ScanCmd.MarkFlagRequired("target")

	ScanCmd.Flags().StringVarP(&options.Ports, "ports", "p", "", 
`Define which ports to scan on the target. You can specify a single port, a range of ports (e.g., 1-1024),
or multiple ports separated by commas.
`)

	ScanCmd.Flags().StringVarP(&options.Script, "script", "s", "",
`Indicate the Nmap scripts to be used during the scan.
You can provide a specific script or a category of scripts (e.g., default, vuln).
`)

	ScanCmd.Flags().StringVarP(&options.Timing, "timing", "", "",
`Set the timing template for the scan, which controls the speed
and stealthiness of the scan. Timing templates range from 0 
(paranoid) to 5 (insane), with higher numbers increasing speed
but also the risk of detection.
`)

	ScanCmd.Flags().BoolVar(&options.SkipHostDiscovery, "skip-host-discovery",false,
`Disable host discovery phase, which skips the process of
identifying live hosts and assumes that the targets are up.
`)

	ScanCmd.Flags().BoolVarP(&options.Aggressive, "aggressive", "a", false,
`Enable aggressive scan options, which include service detection, 
OS detection, and additional scanning techniques to gather more information.
`)

	ScanCmd.Flags().BoolVar(&options.ServiceDetection, "service-detection", false,
`Enable detection of services running on the open 
ports of the target.This helps in identifying the applications
and versions running on the target.
`)

	ScanCmd.Flags().BoolVar(&options.OsDetection, "os-detection", false,
`Enable OS detection to determine the operating system
of the target host. This uses various techniques to guess the 
OS based on network responses.
`)

	ScanCmd.Flags().StringVar(&options.ScanType, "scan-type", "", 
`Specify scan type (e.g., sS for SYN, sT for TCP connect, sU for UDP)
`)
}
