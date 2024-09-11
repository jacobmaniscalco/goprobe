package scan
import (
	"github.com/Ullaakut/nmap/v3"
)

func WithTargets(target string) nmap.Option {
	return nmap.WithTargets(target)
}

func WithPortRange(ports string) nmap.Option {
	return nmap.WithPorts(ports)

}

func WithServiceInfo() nmap.Option {
	return nmap.WithServiceInfo()
}

func WithVerbose(level int) nmap.Option {
	return nmap.WithVerbosity(level)
}

func WithDefaultScript() nmap.Option {
	return nmap.WithDefaultScript()
}

func WithOSDiscovery() nmap.Option {
	return nmap.WithOSDetection()
}
