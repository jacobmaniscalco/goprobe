package ble

import (
	"fmt"
	"github.com/jacobmaniscalco/goprobe/internal/utils/nrfutil"
)

func StartSniffer(macAddress string) error {
	 err := nrfutil.ReadSerial(macAddress)

	 if err != nil {
		 return fmt.Errorf("ReadSerial Error: %v", err)
	 }
	 return nil
}
