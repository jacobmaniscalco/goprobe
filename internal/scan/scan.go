package scan

import (
	"context"
	"fmt"
	"strings"
	"errors"

	"github.com/Ullaakut/nmap/v3"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jacobmaniscalco/blue-caterpillar-cli/ui/styles"
)

type scanDoneMsg struct {
	result string
	err    error
}

type ScanOptions struct {
	 Target string 
	 Ports string
	 Verbose bool
	 Script string
	 Timing string
	 SkipHostDiscovery bool
	 Aggressive bool
	 ServiceDetection bool
	 OsDetection bool
	 ScanType string
}

func PerformScan(options ScanOptions) (string, error) {

	nmapOptions := []nmap.Option{
		nmap.WithTargets(options.Target),
	}

	if options.Ports != "" {
		nmapOptions = append(nmapOptions, nmap.WithPorts(options.Ports))
	}

	if options.Timing != "" {

		if timing, err := mapTiming(options.Timing); err != nil {
			return "", fmt.Errorf("Error in flags: %v", err)

		} else {
			nmapOptions = append(nmapOptions, nmap.WithTimingTemplate(timing))
		}
	}

	if options.SkipHostDiscovery == true {
		nmapOptions = append(nmapOptions, nmap.WithSkipHostDiscovery())
	}

	if options.Aggressive == true {
		nmapOptions = append(nmapOptions, nmap.WithAggressiveScan())
	}

	if options.ServiceDetection == true {
		nmapOptions = append(nmapOptions, nmap.WithServiceInfo())
	}

	if options.OsDetection == true {
		nmapOptions = append(nmapOptions, nmap.WithOSDetection())
	}

	if options.ScanType != "" {
		
		switch options.ScanType {

		case "sS":
			nmapOptions = append(nmapOptions, nmap.WithSYNScan())

		case "sU":
			nmapOptions = append(nmapOptions, nmap.WithUDPScan())
		}

	}

	if options.Script != "" {
		nmapOptions = append(nmapOptions, nmap.WithScripts(options.Script))
	}


	scanner, err := nmap.NewScanner(context.Background(), nmapOptions...)
	if err != nil {
		return "", fmt.Errorf("Failed to create scanner: %v", err)
	}

	results, warnings, err := scanner.Run()
	if err != nil {
		return "", fmt.Errorf("Unable to run nmap scan: %v", err)
	}

	formatted_results := formatScanResults(results, warnings)
	return formatted_results, nil
}

// For some reason, bubbletea and strings builder have an issue with formatting new line characters
// Because of this, some new line characters needs to be in their own statement
func formatScanResults(results *nmap.Run, warnings *[]string) string {

	var sb strings.Builder

	if len(*warnings) > 0 {
		sb.WriteString(styles.WarningStyle.Render(
			fmt.Sprintf("run finished with warnings: %s\n", *warnings)))
	}

	sb.WriteString(styles.TitleStyle.Render("Scan Results"))
	sb.WriteString("\n")

	for _, host := range results.Hosts {
		sb.WriteString(styles.IPStyle.Render(
			fmt.Sprintf("Host: %s", host.Addresses[0].Addr)))
		sb.WriteString("\n")

		for _, match := range host.OS.Matches {
			sb.WriteString(styles.OSStyle.Render(
				fmt.Sprintf("OS Detected: %s %d%%", match.Name, match.Accuracy)))
		}
		sb.WriteString("\n")

		var rows [][]string
		for _, port := range host.Ports {
			rows = append(rows, []string{
				fmt.Sprintf("%d", port.ID),
				fmt.Sprintf("%s", port.State),
				fmt.Sprintf("%s", port.Service),
			})
		}

		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == 0:
					return styles.HeaderStyle
				case row%2 == 0:
					return styles.EvenRowStyle
				default:
					return styles.OddRowStyle
				}
			}).
			Headers("Port", "State", "Service").
			Rows(rows...)

		sb.WriteString(t.Render())
	}
	return sb.String()
}

func mapTiming(timing string) (nmap.Timing, error) {

	switch timing {
	case "T0":
		return nmap.TimingSlowest, nil 
	case "T1":
		return nmap.TimingSneaky, nil
	case "T2":
		return nmap.TimingPolite, nil
	case "T3":
		return nmap.TimingNormal, nil
	case "T4": 
		return nmap.TimingAggressive, nil
	case "T5":
		return nmap.TimingFastest, nil
	default:
		return 0, errors.New("Invalid timing.")
	}

}
