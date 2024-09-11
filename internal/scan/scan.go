package scan

import (
	"fmt"
	"strings"
	"context"

	"github.com/jacobmaniscalco/blue-caterpillar-cli/ui/styles"
	"github.com/Ullaakut/nmap/v3"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func RunScan(target string, ports string) (*nmap.Run, error) {
	
	
		
	options := []nmap.Option {			
		WithTargets(target),
		WithServiceInfo(),
		WithOSDiscovery(),
		nmap.WithScripts("vuln"),
	}

	if ports != "" {
		options = append(options, WithPortRange(ports))
	}

	scanner, err := nmap.NewScanner(context.Background(), options...)
	if err != nil {			
		return nil, fmt.Errorf("Failed to create scanner: %v", err)
	}
	
	results, warnings, err := scanner.Run()
	if err != nil {
		return nil, fmt.Errorf("Unable to run nmap scan: %v", err)
	}

	
	printScanResults(results, warnings)
	return results, err
}


func printScanResults(results *nmap.Run, warnings *[]string) {

	var sb strings.Builder 
	
	if len(*warnings) > 0 {
		sb.WriteString(styles.WarningStyle.Render(fmt.Sprintf("run finished with warnings: %s\n", *warnings)))
	}

	sb.WriteString(styles.TitleStyle.Render("Scan Results"))
	sb.WriteString("\n")
	for _, host := range results.Hosts {

		sb.WriteString(styles.IPStyle.Render(fmt.Sprintf("Host: %s", host.Addresses[0].Addr)))

		// for some reason, this newline character has to be on a seperate line to avoid formatting issues
		sb.WriteString("\n")

		for _, match := range host.OS.Matches {
			sb.WriteString(styles.OSStyle.Render(fmt.Sprintf("OS Detected: %s %d%%", match.Name, match.Accuracy)))
		}
		
		sb.WriteString("\n")
		var rows [][]string
		for _, port := range host.Ports {
			rows = append(rows, []string {
				fmt.Sprintf("%d", port.ID),
				fmt.Sprintf("%s", port.State),
				fmt.Sprintf("%s", port.Service),
			})
		}

		t := table.New().
			   Border(lipgloss.NormalBorder()).
			   BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			   StyleFunc(func (row, col int) lipgloss.Style {
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

	fmt.Println(sb.String())
}

