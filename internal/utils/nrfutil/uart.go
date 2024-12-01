package nrfutil

import (
	"encoding/binary"
	"fmt"
	"github.com/tarm/serial"
	"log"
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

	PROTOVER_V1           = 1
	PROTOVER_V2           = 2
	PROTOVER_V3           = 3
	EVENT_PACKET_ADV_PDU  = 0x02
	EVENT_FOLLOW = 0x01
	EVENT_PACKET_DATA_PDU = 0x06
	PHY_CODED             = 2
	PING_RESP             = 0x0E
	RESP_VERSION          = 0x1C
	RESP_TIMESTAMP        = 0x1E
	SWITCH_BAUD_RATE_RESP = 0x14
	SWITCH_BAUD_RATE_REQ  = 0x13

	BLE_HEADER_LENGTH  = 10

	PAYLOAD_LEN_POS_V1 = 1
	PAYLOAD_LEN_POS    = 0
	PROTOVER_POS       = PAYLOAD_LEN_POS + 2
	PACKETCOUNTER_POS  = PROTOVER_POS + 1
	ID_POS             = PACKETCOUNTER_POS + 2
	BLE_HEADER_LEN_POS = ID_POS + 1
	FLAGS_POS          = BLE_HEADER_LEN_POS + 1
	CHANNEL_POS        = FLAGS_POS + 1
	RSSI_POS           = CHANNEL_POS + 1
	EVENTCOUNTER_POS   = RSSI_POS + 1
	TIMESTAMP_POS      = EVENTCOUNTER_POS + 2
	BLEPACKET_POS      = TIMESTAMP_POS + 4
	PAYLOAD_POS        = BLE_HEADER_LEN_POS

	PACKET_TYPE_DATA        = 0x02
	PACKET_TYPE_ADVERTISING = 0x01
)

type Packet struct {
	packetList      []uint8
	version         string
	protover        uint8
	packetCounter   uint16
	payloadLength   uint16
	id              uint8
	OK              bool
	crcOK           bool
	valid           bool
	bleHeaderLength uint8
	flags           uint8
	channel         uint8
	rawRSSI         uint8
	RSSI            int8
	phy             uint8
	eventCounter    uint16
	timestamp       uint32
	direction       bool
	encrypted       bool
	micOK           bool
	baudRate        uint32
	blePacket       *BLEPacket
}

type BLEPacket struct {
	packetList []uint8
	accessAddress uint32
	packetType      uint8
	length uint8
	payload         []uint8
	coded           bool
	codingIndicator uint8
	advType         uint8
	txAddrType      uint8
	rxAddrType      uint8
	llid            uint8
	sn              uint8
	nesn            uint8
	md              uint8
	advAddress      []uint8
	scanAddress     []uint8
	name string
}

type packetReader struct {
	packetCounter int
}

var PORT *serial.Port
var MAC_ADDRESS string 

func ReadSerial(macAddress string) error {
	MAC_ADDRESS = macAddress
	port, err := serial.OpenPort(&serial.Config{
		Name:        "/dev/ttyACM0",
		Baud:        1000000,
		ReadTimeout: time.Second * 1,
	})
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("ReadSerial Error: %v", err)
	}
	PORT = port
	defer port.Close()

	for {
		packetList, err := decodeFromSLIP(port)
		if err != nil {
			fmt.Printf("Error decoding SLIP packet: %v\n", err)
			continue
		}
		fmt.Printf("Raw Packet Data(Hex): %X\n", packetList)
		decodeSnifferPacket(packetList)
	}
}

func toLittleEndian(value uint16, length int) []byte {
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
		}
		if time.Since(timeStart) > timeout {
			return nil, fmt.Errorf("timeout waiting for SLIP_START")
		}
	}
	for !endOfPacket {
		serialByte, err := getSerialByte(port)
		if err != nil {
			return nil, fmt.Errorf("Failed during packet decoding: %w", err)
		}
		switch serialByte {
		case SLIP_END:
			endOfPacket = true
		case SLIP_ESC:
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

func decodeSnifferPacket(packetList []byte) {

	protocolNumber := packetList[PROTOVER_POS]
	packetCounter := binary.LittleEndian.Uint16(packetList[PACKETCOUNTER_POS : PACKETCOUNTER_POS+2])
	payloadLength := binary.LittleEndian.Uint16(packetList[PAYLOAD_LEN_POS : PAYLOAD_LEN_POS+2])
	id := packetList[ID_POS]

	packet := &Packet{
		packetList:    packetList,
		protover:      protocolNumber,
		packetCounter: packetCounter,
		payloadLength: payloadLength,
		id:            id,
	}
	fmt.Printf("Payload Length (Hex): % X\n", packet.payloadLength)
	fmt.Printf("Payload Length (Dec): %d\n", packet.payloadLength)
	fmt.Printf("Protocol Version (Hex): % X\n", packet.protover)
	fmt.Printf("Protocol Version (Dec): %d\n", packet.protover)
	fmt.Printf("Packet Counter (Hex): % X\n", packet.packetCounter)
	fmt.Printf("Packet Counter (Dec): %d\n", packet.packetCounter)
	fmt.Printf("Packet ID (Hex): % X\n", packet.id)
	fmt.Printf("Packet ID (Dec): %d\n", packet.id)

	readPayload(packet)

	if packet.OK {
		fmt.Printf("Payload Length (Hex): % X\n", packet.bleHeaderLength)
		fmt.Printf("payload length (Dec): %d\n", packet.bleHeaderLength)
		
		fmt.Printf("Flags (Hex): % X\n", packet.flags)
		fmt.Printf("Flags (Dec): %d\n", packet.flags)

		fmt.Printf("Channel Index (Hex): % X\n", packet.channel)
		fmt.Printf("Channel Index (Dec): %d\n", packet.channel)
		
		fmt.Printf("RSSI (Hex): % X\n", packet.RSSI)
		fmt.Printf("RSSI (Dec): %d\n", packet.RSSI)
		
		fmt.Printf("Event Counter (Hex): % X\n", packet.eventCounter)
		fmt.Printf("Event Counter (Dec): %d\n", packet.eventCounter)
		
		fmt.Printf("Timestamp (Hex): % X\n", packet.timestamp)
		fmt.Printf("Timestamp (Dec): %d\n", packet.timestamp)

		fmt.Printf("Access Address (Hex): % X\n", packet.blePacket.accessAddress)
		fmt.Printf("Access Address (Dec): %d\n", packet.blePacket.accessAddress)
		
		//fmt.Printf("Advertising Address (Hex): % X\n", packet.blePacket.advAddress)
		//fmt.Printf("Advertising Address (Dec): %d\n", packet.blePacket.advAddress)
		macBytes := packet.blePacket.advAddress[:6]

		mac := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", macBytes[0], macBytes[1], macBytes[2], macBytes[3], macBytes[4], macBytes[5])


		fmt.Println("Advertising Address length: ", len(packet.blePacket.advAddress))
		
		fmt.Printf("BLE Packet Type (Hex): % X\n", packet.blePacket.packetType)
		fmt.Printf("BLE Packet Type (Dec): %d\n", packet.blePacket.packetType)
		
		fmt.Printf("BLE Packet Scan Address (Hex): % X\n", packet.blePacket.scanAddress)
		fmt.Printf("ble Packet Scan Address (Dec): %d\n", packet.blePacket.scanAddress)
		
		fmt.Println("BLE Packet Payload: ", packet.blePacket.payload)
		fmt.Println("BLE Packet Payload Length: ", len(packet.blePacket.payload))


		reader := &packetReader{
			packetCounter: 0,
		}
		fmt.Println("mac: ", mac)
		if mac == MAC_ADDRESS {
			fmt.Println("FOUND RASPBERRY PI!")
			err := SendFollow(macBytes, false, false, false, reader)

			if err != nil {
				fmt.Println("ERROR: ", err)
			}
		}
	}

}

func readFlags(packet *Packet) {
	packet.crcOK = (packet.flags & 1) != 0
	packet.direction = (packet.flags & 2) != 0
	packet.encrypted = (packet.flags & 4) != 0
	packet.micOK = (packet.flags & 8) != 0
	packet.phy = (packet.flags >> 4) & 7
	packet.OK = packet.crcOK && (packet.micOK || !packet.encrypted)

}

func validatePacketList(packet *Packet) bool {
	if packet.payloadLength + HEADER_LENGTH == uint16(len(packet.packetList)) {
		return true
	} else {
		return false
	}
}


func readPayload(packet *Packet) {

	if !validatePacketList(packet) {
		fmt.Println("Packet not validated.")
	}

	switch packet.id {
	case EVENT_PACKET_ADV_PDU, EVENT_PACKET_DATA_PDU:
		packet.bleHeaderLength = packet.packetList[BLE_HEADER_LEN_POS]
		if packet.bleHeaderLength == BLE_HEADER_LENGTH {
			packet.flags = packet.packetList[FLAGS_POS]
			readFlags(packet)
			packet.channel = packet.packetList[CHANNEL_POS]
			packet.rawRSSI = packet.packetList[RSSI_POS]
			packet.RSSI = int8(packet.rawRSSI) * -1
			packet.eventCounter = binary.LittleEndian.Uint16(packet.packetList[EVENTCOUNTER_POS : EVENTCOUNTER_POS+2])

			packet.timestamp = binary.LittleEndian.Uint32(packet.packetList[TIMESTAMP_POS : TIMESTAMP_POS+4])

			// removing a padding byte and update payload length in the packet list
			if packet.phy == PHY_CODED {
				index := BLEPACKET_POS + 6 + 1
				packet.packetList = append(packet.packetList[:index], packet.packetList[index+1:]...)
				fmt.Println("Inside PHY_CODED, index at: ", index)

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
			} 
			

		} else {
			fmt.Printf("Invalid BLE Header Length: %d\n", packet.bleHeaderLength)
			packet.valid = false
		}
		
		if packet.OK {
			var packetType uint8
			if packet.protover >= PROTOVER_V3 {
				if packet.id == EVENT_PACKET_ADV_PDU {
					packetType = PACKET_TYPE_ADVERTISING
				} else {
					packetType = PACKET_TYPE_DATA
				}
			}

			packet.blePacket = &BLEPacket{
				packetType: packetType,
				packetList:    packet.packetList[BLEPACKET_POS:],
			}

			decodeBLEPacket(packet)
			parseBLEPayload(packet.blePacket.payload)
		}

	case EVENT_FOLLOW:
		fmt.Println("EVENT FOLLOW")
	case PING_RESP:
		if packet.protover < PROTOVER_V3 {
			fmt.Printf("Protocol under v3")
			//packet.version = parseLittleEndian(packet.packetList[PAYLOAD_POS:PAYLOAD_POS+2])
		}
	case RESP_VERSION:
		packet.version = string(packet.packetList[PAYLOAD_POS:])

	case RESP_TIMESTAMP:
		fmt.Println("Using RESP TIMESTAMP")
		packet.timestamp = binary.LittleEndian.Uint32(packet.packetList[PAYLOAD_POS : PAYLOAD_POS+4])

	case SWITCH_BAUD_RATE_RESP, SWITCH_BAUD_RATE_REQ:
		packet.baudRate = binary.LittleEndian.Uint32(packet.packetList[PAYLOAD_POS : PAYLOAD_POS+4])

	}
}

func decodeBLEPacket(packet *Packet) {
	var offset uint8 = 0
	offset = extractAccessAddress(packet, offset)
	offset = extractFormat(packet, offset)

	if packet.blePacket.packetType == PACKET_TYPE_ADVERTISING {
		offset = extractAdvHeader(packet, offset)
	} else {
		offset = extractConnHeader(packet, offset)
	}
	offset = extractLength(packet, offset)

	packet.blePacket.payload = packet.blePacket.packetList[offset:]

	if packet.blePacket.packetType == PACKET_TYPE_ADVERTISING {
		offset = extractAddresses(packet, offset)
		extractName(packet, offset)
	}

}

func extractName(packet *Packet, offset uint8) {
	name := ""
	validTypes := []uint8{0, 2, 4, 6, 7} // List of valid advTypes

	if contains(packet.blePacket.advType, validTypes) {
		i := int(offset)
		for i < len(packet.packetList) {
			length := int(packet.packetList[i])
			if i+length+1 > len(packet.packetList) || length == 0 {
				break
			}
			packetType := packet.packetList[i+1]
			if packetType == 8 || packetType == 9 {
				nameList := packet.packetList[i+2 : i+length+1]
				name = ""
				for _, b := range nameList {
					name += string(b)
				}
			}
			i += length + 1
		}
		name = `"` + name + `"`
	} else if packet.blePacket.advType == 1 {
		name = "[ADV_DIRECT_IND]"
	}

	packet.blePacket.name = name
}


func extractLength(packet *Packet, offset uint8) uint8 {
	packet.blePacket.length = packet.blePacket.packetList[offset]
	return offset + 1
}

func extractAccessAddress(packet *Packet, offset uint8) uint8 {
	packet.blePacket.accessAddress = binary.LittleEndian.Uint32(packet.blePacket.packetList[offset : offset+4])
	return offset + 4
}

func extractFormat(packet *Packet, offset uint8) uint8 {
	packet.blePacket.coded = packet.phy == PHY_CODED
	if packet.phy == PHY_CODED {
		packet.blePacket.coded = true
	} else {
		packet.blePacket.coded = false
	}
	if packet.blePacket.coded {
		packet.blePacket.codingIndicator = packet.blePacket.packetList[offset] & 3
		return offset + 1
	}
	return offset
}

func extractAdvHeader(packet *Packet, offset uint8) uint8 {
	packet.blePacket.advType = packet.blePacket.packetList[offset] & 15
	packet.blePacket.txAddrType = (packet.blePacket.packetList[offset] >> 6) & 1
	if packet.blePacket.advType == 1 || packet.blePacket.advType == 3 || packet.blePacket.advType == 5 {
		packet.blePacket.rxAddrType = (packet.blePacket.packetList[offset] << 7) & 1
	} else if packet.blePacket.advType == 7 {
		flags := packet.blePacket.packetList[offset+2]
		if flags&0x02 != 0 {
			packet.blePacket.rxAddrType = (packet.blePacket.packetList[offset] << 7) & 1
		}
	}
	return offset + 1
}

func extractConnHeader(packet *Packet, offset uint8) uint8 {
	packet.blePacket.llid = packet.blePacket.packetList[offset] & 3
	packet.blePacket.sn = (packet.blePacket.packetList[offset>>2]) & 1
	packet.blePacket.nesn = (packet.blePacket.packetList[offset] >> 3) & 1
	packet.blePacket.md = (packet.blePacket.packetList[offset] >> 4) & 1
	return offset + 1
}

func reverseBytes(input []byte) []byte {
    output := make([]byte, len(input))
    for i, v := range input {
        output[len(input)-1-i] = v
    }
    return output
}

func extractAddresses(packet *Packet, offset uint8) uint8 {
	var addr []byte
	var scanAddr []byte

	validTypes := []byte{0, 1, 2, 4, 6}
	if contains(packet.blePacket.advType, validTypes) {
		addr = reverseBytes(packet.blePacket.packetList[offset : offset+6])
		addr = append(addr, packet.blePacket.txAddrType)
		offset += 6
	}

	if packet.blePacket.advType == 3 || packet.blePacket.advType == 5 {
		scanAddr = reverseBytes(packet.blePacket.packetList[offset:offset+6])
		scanAddr = append(scanAddr, packet.blePacket.txAddrType)
		offset += 6
		addr = reverseBytes(packet.blePacket.packetList[offset:offset+6])
		addr = append(addr, packet.blePacket.rxAddrType)
		offset += 6
	}

	if packet.blePacket.advType == 1 {
		scanAddr = reverseBytes(packet.blePacket.packetList[offset:offset+6])
		scanAddr = append(scanAddr, packet.blePacket.rxAddrType)
		offset += 6
	}

	if packet.blePacket.advType == 7 {
		ext_header_len := packet.blePacket.packetList[offset] & 0x3F
		offset += 1

		ext_header_offset := offset
		flags := packet.blePacket.payload[offset]
		ext_header_offset += 1

		if flags&0x01 != 0 {
			addr = reverseBytes(packet.blePacket.packetList[ext_header_offset:ext_header_offset+6])
			addr = append(addr, packet.blePacket.txAddrType)
			ext_header_offset += 6
		}

		if flags&0x02 != 0 {
			scanAddr = reverseBytes(packet.blePacket.packetList[ext_header_offset:ext_header_offset+6])
			scanAddr = append(scanAddr, packet.blePacket.rxAddrType)
			ext_header_offset += 6
		}
		offset += ext_header_len
	}

	packet.blePacket.advAddress = addr
	packet.blePacket.scanAddress = scanAddr
	return offset
}

func reverse(slice []byte) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func contains(value uint8, slice []byte) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}
	return false
}

func parseLittleEndian(packet []byte) uint32 {

	var total uint32 = 0
	for i := 0; i < len(packet); i++ {
		total += uint32(packet[i]) << (8 * i)
	}

	return total
}

func parseTimeStamp(packet []byte) uint32 {
	return binary.LittleEndian.Uint32(packet)
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

func formatAddress(address uint32) string {
// Convert uint32 to 4 bytes
	b1 := byte(address >> 24) // Most significant byte
	b2 := byte(address >> 16)
	b3 := byte(address >> 8)
	b4 := byte(address)

	// Add padding to make it 6 bytes (00:00 + input bytes)
	mac := []byte{b1, b2, b3, b4}

	// Format as a MAC address string
	return fmt.Sprintf("%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3])
}

func sendPacket(id byte, payload []uint8, packetReader *packetReader) error {
	packet := []byte{HEADER_LENGTH, byte(len(payload)), PROTOVER_V1}
	packet = append(packet, toLittleEndian(7, 2)...)
	packet = append(packet, id)
	packet = append(packet, payload...)
	packetReader.packetCounter++
	packet = encodeToSLIP(packet)
	fmt.Println("Inside sendPacket, sending: \n", packet)
	return writeList(packet)

}

func SendFollow(addr []byte, followOnlyAdvertisements bool, followOnlyLegacy bool, followCoded bool, packetReader *packetReader) error {
	flags0 := byte(0)
	if followOnlyAdvertisements {
		flags0 |= 1
	}
	if followOnlyLegacy {
		flags0 |= 1 << 1
	}
	if followCoded {
		flags0 |= 1 << 2
	}
	fmt.Printf("Follow flags: %08b", flags0)
	addr = append(addr, byte(0))
	return sendPacket(REQ_FOLLOW, append(addr, flags0), packetReader)
}

func writeList(array []byte) error {
	//var foo = []byte {171, 6, 16, 1, 5, 0, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 188}
	fmt.Println("Sending writeList: ", array)
	//array = foo
	n, err := PORT.Write(array)

	if err != nil {
		fmt.Println("Error writing to port.")
		return err
	}

	fmt.Println("Output from write: ", n)
	return nil
}
