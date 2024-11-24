package nrfutil

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"strings"
	"time"
)

const (
	SLIP_START     = 0xAB
	SLIP_END       = 0xBC
	SLIP_ESC       = 0xCD
	SLIP_ESC_START = SLIP_START + 1
	SLIP_ESC_END   = SLIP_END + 1
	SLIP_ESC_ESC   = SLIP_ESC + 1
	REQ_FOLLOW     = 0x00
	HEADER_LENGTH  = 6
	
	PROTOVER_V1    = 1
	PROTOVER_V2 = 2
	PROTOVER_V3 = 3
	EVENT_PACKET_ADV_PDU = 0x02
	EVENT_PACKET_DATA_PDU = 0x06
	PHY_CODED = 2

	BLE_HEADER_PEN_POS = ID_POS + 1
	BLE_HEADER_LENGTH = 10

	PAYLOAD_LEN_POS_V1 = 1
	PAYLOAD_LEN_POS = 0
	PROTOVER_POS = PAYLOAD_LEN_POS+2
	PACKETCOUNTER_POS = PROTOVER_POS+1
	ID_POS = PACKETCOUNTER_POS+2
	BLE_HEADER_LEN_POS = ID_POS+1
	FLAGS_POS = BLE_HEADER_LEN_POS+1
	CHANNEL_POS = FLAGS_POS+1
	RSSI_POS = CHANNEL_POS+1
	EVENTCOUNTER_POS = RSSI_POS+1
	TIMESTAMP_POS = EVENTCOUNTER_POS+2
	BLEPACKET_POS = TIMESTAMP_POS+4
	PAYLOAD_POS = BLE_HEADER_LEN_POS

	PACKET_TYPE_DATA = 0x02
	PACKET_TYPE_ADVERTISING = 0x01
)

type Packet struct {
	packetList []byte
	protover uint16
	packetCounter uint8
	payloadLength uint8
	id uint8
	OK bool
	crcOK bool
	valid bool
	bleHeaderLength uint8
	flags uint8
	channel uint8
	rawRSSI uint8
	RSSI uint8
	phy uint8
	eventCounter uint8
	timestamp uint8
	direction bool
	encrypted bool
	micOK bool
	blePacket BLEPacket 
}

type BLEPacket struct {

	pType uint8 
	payload []byte
	accessAddress []byte

}

func ReadSerial(macAddress string) error {
	port, err := serial.OpenPort(&serial.Config{
		Name:        "/dev/ttyACM0",
		Baud:        1000000,
		ReadTimeout: time.Second * 1,
	})
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("ReadSerial Error: %v", err)
	}
	defer port.Close()

	for {
		packetList, err := decodeFromSLIP(port)
		if err != nil {
			fmt.Printf("Error decoding SLIP packet: %v\n", err)
			continue
		}
		fmt.Printf("Raw Packet Data: %X\n", packetList)
		fmt.Printf("Packet length: %d\n", len(packetList))
		decodeSnifferPacket(packetList)
	}
}

func toLittleEndian(value uint8, length int) []byte {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(value >> (8 * i) & 0xFF)
	}
	return bytes
}

func getSerialByte(port *serial.Port) (byte, error) {
	buf := make([]byte, 1)
	n, err := port.Read(buf)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, fmt.Errorf("timeout: no data received")
	}
	return buf[0], nil
}

func decodeFromSLIP(port *serial.Port) ([]byte, error) {

	var dataBuffer []byte
	startOfPacket := false
	endOfPacket := false

	timeout := time.Second * 5
	timeStart := time.Now()

	for !startOfPacket {
		res, err := getSerialByte(port)
		if err != nil {
			return nil, fmt.Errorf("failed to find SLIP_START: %w", err)
		}
		if res == SLIP_START {
			startOfPacket = true
			fmt.Println("SLIP_START found, packet starting")
		}
		if time.Since(timeStart) > timeout {
			return nil, fmt.Errorf("timeout waiting for SLIP_START")
		}
	}
	for !endOfPacket {
		serialByte, err := getSerialByte(port)
		if err != nil {
			fmt.Printf("Error reading byte: %v", err)
			return nil, fmt.Errorf("Failed during packet decoding: %w", err)
		}
		switch serialByte {
		case SLIP_END:
			endOfPacket = true
			fmt.Println("SLIP_END found, packet complete")
		case SLIP_ESC:
			fmt.Printf("SLIP_ESC sequence found: 0x%X\n", serialByte)
			serialByte, err = getSerialByte(port)
			if err != nil {
				fmt.Printf("Error after SLIP_ESC: %v", err)
				return nil, fmt.Errorf("Failed to read after SLIP_ESC: %w", err)
			}
			//fmt.Printf("SLIP_ESC sequence found: 0x%X", nextByte)
			switch serialByte {
			case SLIP_ESC_START:
				dataBuffer = append(dataBuffer, SLIP_START)
			case SLIP_ESC_END:
				dataBuffer = append(dataBuffer, SLIP_END)
			case SLIP_ESC_ESC:
				dataBuffer = append(dataBuffer, SLIP_ESC)
			default:
				dataBuffer = append(dataBuffer, SLIP_END)
				fmt.Printf("Unknown SLIP_ESC sequence: 0x%X\n", serialByte)
				return nil, fmt.Errorf("Invalid escape sequence: 0x%X", serialByte)
			}
		default:
			dataBuffer = append(dataBuffer, serialByte)
		}

		if time.Since(timeStart) > timeout {
			return nil, fmt.Errorf("timeout while waiting for SLIP_END")
		}
	}
	
	return dataBuffer, nil
}

func encodeToSLIP(data []byte) []byte {

	encoded := []byte{SLIP_START}
	for _, b := range data {
		switch b {
		case SLIP_START:
			encoded = append(encoded, SLIP_ESC, SLIP_ESC_START)
		case SLIP_END:
			encoded = append(encoded, SLIP_ESC, SLIP_ESC_END)
		case SLIP_ESC:
			encoded = append(encoded, SLIP_ESC, SLIP_ESC_ESC)
		default:
			encoded = append(encoded, b)
		}
	}
	encoded = append(encoded, SLIP_END)
	return encoded
}

func macAddressToBytes(mac string) ([]byte, error) {

	cleaned := strings.ReplaceAll(mac, ":", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	bytes, err := hex.DecodeString(cleaned)
	if err != nil {
		return nil, fmt.Errorf("Invalid MAC address: %w", err)
	}

	if len(bytes) != 6 {
		return nil, fmt.Errorf("MAC address must be 6 bytes, got %d bytes.", len(bytes))
	}

	return bytes, nil
}

func decodeSnifferPacket(packetList []byte) {

	protocolNumber := packetList[PROTOVER_POS]
	packetCounter := parseLittleEndian(packetList[PACKETCOUNTER_POS:PACKETCOUNTER_POS+2])
	payloadLength := parseLittleEndian(packetList[PAYLOAD_LEN_POS:PAYLOAD_LEN_POS+2])
	id := packetList[ID_POS]

	packet := &Packet {
		packetList: packetList,
		protover: uint16(protocolNumber),
		packetCounter: packetCounter,
		payloadLength: payloadLength,
		id: id,
	}
	fmt.Printf("Packet Counter: %d\n", packetCounter)
	fmt.Printf("Protocol Version: %d\n", protocolNumber)
	fmt.Printf("Payload Length: %d\n", payloadLength)

	readPayload(packet.packetList[PAYLOAD_POS: PAYLOAD_POS + packet.payloadLength], packet)

}

func readFlags(packet *Packet) {
	packet.crcOK = (packet.flags & 1) != 0
	packet.direction = (packet.flags & 2) != 0 
	packet.encrypted = (packet.flags & 4) != 0 
	packet.micOK = (packet.flags & 8) != 0
	packet.phy =  (packet.flags >> 4) & 7
	packet.OK =  packet.crcOK && (packet.micOK || !packet.encrypted)

}

func readPayload(payload []byte, packet *Packet) {

	switch packet.id {
	case EVENT_PACKET_ADV_PDU, EVENT_PACKET_DATA_PDU:
		// get ble header, flags, and more
		packet.bleHeaderLength = payload[BLE_HEADER_LEN_POS]
		if packet.bleHeaderLength == BLE_HEADER_LENGTH {
			packet.flags = payload[FLAGS_POS]
			readFlags(packet)
			packet.channel = payload[CHANNEL_POS]
			packet.rawRSSI = payload[RSSI_POS]
			packet.RSSI = - packet.rawRSSI
			packet.eventCounter = parseLittleEndian(payload[EVENTCOUNTER_POS: EVENTCOUNTER_POS+2])
			packet.timestamp = parseLittleEndian(payload[TIMESTAMP_POS: TIMESTAMP_POS+4])

			// removing a padding byte and update payload length in the packet list
			if packet.phy == PHY_CODED {
				index := BLEPACKET_POS + 6 + 1
				packet.packetList = append(packet.packetList[:index], packet.packetList[index+1:]...)	
			} else {
				index := BLEPACKET_POS + 6
				packet.packetList = append(packet.packetList[:index], packet.packetList[index+1:]...)
			}
			// update length accordingly
			packet.payloadLength -= 1

			if packet.protover >= PROTOVER_V2 {
				payloadLength := toLittleEndian(packet.payloadLength, 2)
				packet.packetList[PAYLOAD_LEN_POS] = payloadLength[0]
				packet.packetList[PAYLOAD_LEN_POS+1] = payloadLength[1]
			} else {
				packet.packetList[PAYLOAD_LEN_POS_V1] = packet.payloadLength
			}

		} else {
			fmt.Errorf("Invalid BLE Header Length: %v", packet.bleHeaderLength)
			packet.valid = false
		}
		
		if packet.OK {
			if packet.protover >= PROTOVER_V3 {
				if packet.id == EVENT_PACKET_ADV_PDU {
					packetType := PACKET_TYPE_ADVERTISING
				} else {
					packetType := PACKET_TYPE_DATA
				}
			}

			packet.blePacket = blePacket(
		}

	case PING_RESP:
	
	case RESP_VERSION:
	
	case RESP_TIMESTAMP:

	case SWITCH_BAUD_RATE_RESP, SWITCH_BAUD_RATE_REQ:

	}

}

func decodeBLEPacket(packet []byte) {
	if len(packet) < 7 {
		fmt.Printf("Packet too short to decode (length %d): %X\n", len(packet), packet)
		return
	}
	preamble := packet[0]
	accessAddress := packet[0:4]
	header := packet[4:7]
	payload := packet[6:]

	// Print extracted fields
	fmt.Printf("Preamble: 0x%02X\n", preamble)
	accessAddressStr := fmt.Sprintf("%02X:%02X:%02X:%02X", accessAddress[0], accessAddress[1], accessAddress[2], accessAddress[3])
	fmt.Printf("Access Address: %s\n", accessAddressStr)
	//fmt.Printf("Access Address: %X\n", accessAddress)
	fmt.Printf("Header: 0x%02X\n", header)
	fmt.Printf("Payload: %X\n", payload)

	pduType := header[0] & 0x0F
	fmt.Println("header: ", header)
	fmt.Println("pduType: ", pduType)

	switch pduType {
    case 0x00:
        fmt.Println("Packet Type: ADV_IND (Connectable undirected advertising)")
    case 0x01:
        fmt.Println("Packet Type: ADV_DIRECT_IND (Connectable directed advertising)")
    case 0x02:
        fmt.Println("Packet Type: ADV_NONCONN_IND (Non-connectable undirected advertising)")
    case 0x03:
        fmt.Println("Packet Type: SCAN_REQ (Scan request)")
    case 0x04:
        fmt.Println("Packet Type: SCAN_RSP (Scan response)")
    case 0x05:
        fmt.Println("Packet Type: CONNECT_REQ (Connection request)")
    case 0x06:
        fmt.Println("Packet Type: ADV_SCAN_IND (Scannable undirected advertising)")
    default:
        fmt.Println("Packet Type: Unknown or Reserved")
    }

	// Parse the payload for basic data
	if len(payload) >= 6 {

		advertisingAddress := payload[:6]
		fmt.Printf("Advertising Address: %02X:%02X:%02X:%02X:%02X:%02X\n",
			advertisingAddress[0], advertisingAddress[1], advertisingAddress[2],
			advertisingAddress[3], advertisingAddress[4], advertisingAddress[5])

		parseBLEPayload(payload[6:])
	} else if len(payload) >= 4 {
		cid := binary.LittleEndian.Uint16(payload[2:4])
		fmt.Printf("CID: 0x%04X\n", cid)
	} else {
		fmt.Println("Unrecognized payload format")
	}
}

func parseLittleEndian(packet []byte) uint8 {

	var total uint8 = 0
	for i := range len(packet) {
		total +=  (packet[i] << (8 * i))
	}

	return total
}

func parseBLEPayload(payload []byte) {
	pos := 0
	for pos < len(payload) {
		if pos >= len(payload) {
			fmt.Println("Payload parsing out of bounds.")
			break
		}

		adLength := int(payload[pos])
		pos++

		if pos+adLength > len(payload) || adLength == 0 {
			fmt.Printf("Invalid length field (%d), at position %d\n", adLength, pos)
			break
		} else {
			fmt.Println("length okay")
		}

		adType := payload[pos]
		adData := payload[pos+1 : pos+adLength]
		pos += adLength

		switch adType {
		case 0x01: // Flags
			fmt.Printf("  Flags: 0x%02X\n", adData[0])
		case 0x02, 0x03:
			for i := 0; i < len(adData); i += 2 {
				uuid := binary.LittleEndian.Uint16(adData[i : i+2])
				fmt.Printf("  Service UUID: 0x%04X\n", uuid)
			}
		case 0x08, 0x09: // Complete Local Name
			fmt.Printf("  Device Name: %s\n", string(adData))
		case 0xFF: // Manufacturer Specific Data
			fmt.Printf("  Manufacturer Data: %X\n", adData)
		default:
			fmt.Printf("  AD Type 0x%02X: %X\n", adType, adData)
		}
	}
}
