package scan

type scanDoneMsg struct {
	result string
	err    error
}

type ScanOptions struct {
	 Host string 
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
