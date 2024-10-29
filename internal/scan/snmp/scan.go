package snmp


import (
	"time"
	"log"
	"fmt"

	"github.com/gosnmp/gosnmp"
	scan "github.com/jacobmaniscalco/goprobe/internal/scan"
)


func ScanSNMPDevice(options scan.ScanOptions, port uint16) {

	params := &gosnmp.GoSNMP{
		Target: options.Host,
		Port: port,
		Community: "public",
		Version: gosnmp.Version2c,
		Timeout: time.Duration(2) * time.Second,
		Retries: 1,
		//Logger: log.New(log.Writer(), "", 0),
	}

	err := params.Connect()
	if err != nil {
		log.Fatalf("Error connecting to target %s: %v", options.Host, err)
	}
	defer params.Conn.Close()

	err = params.Walk(".1.3.6.1.2.1", walkFunc)
	if err != nil {
		log.Printf("SNMP walk failed for target %s: %v", options.Host, err)
	}
}

func walkFunc(pdu gosnmp.SnmpPDU) error {
	switch pdu.Type {
	case gosnmp.OctetString:
		value := string(pdu.Value.([]byte))
		fmt.Printf("OID: %s, Value: %s\n", pdu.Name, value)
	default:
		fmt.Printf("OID: %s, Value: %v\n", pdu.Name, pdu.Value)
	}
	return nil
}

