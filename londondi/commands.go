package londondi

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"strconv"
	"strings"
)

/* to set a state variable, <DI_SETSV>
<node>
<virtual_device>
<object>
<state_variable>
<data>
*/

const LEN_NODE = 2
const LEN_ADDR = 5

func BuildRawMuteCommand(input, address, status string) ([]byte, error) {

	log.Printf("Building raw mute command for input: %s on address: %s", input, address)

	command := []byte{DI_SETSV}
	log.Printf("Command string: %s", hex.EncodeToString(command))

	//the node is the hex representation of the HiQnet Address, which is assumed to be the last 4 digits of the IP address, well... sort of
	firstDigit := strings.Split(address, ".")[1]
	nodeString := firstDigit[len(firstDigit)-1:] + strings.Split(address, ".")[2]

	log.Printf("Detected IP Address: %s", address)
	log.Printf("Detected HiQnet address: %s (decimal), %X (hex)", nodeString, (nodeString))

	nodeBytes, err := hex.DecodeString(nodeString)
	if err != nil {
		errorMessage := "Could not decode node string: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	for {
		if len(nodeBytes) > LEN_NODE {
			return []byte{}, errors.New("detected HiQnet node longer than 2 bytes")
		} else if len(nodeBytes) == LEN_NODE {
			nodeBytes = append([]byte{0x00}, nodeBytes...)
		}
	}

	command = append(command, nodeBytes...)
	log.Printf("Command string: %X", command)

	command = append(command, VIRTUAL_DEVICE)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	gainBlock, err := hex.DecodeString(input)
	if err != nil {
		errorMessage := "Could not decode input string: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	command = append(command, gainBlock...)
	command = append(command, stateVariables["mute"]...)
	command = append(command, muteStates[status]...)

	checksum := GetChecksumByte(command)
	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command, _ = MakeSubstitutions(command, ENCODE)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	stx := []byte{STX}
	command = append(stx, command...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, ETX)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	return command, nil
}

func BuildRawVolumeCommand(input string, address string, volume string) ([]byte, error) {

	log.Printf("Building raw volume command for input: %s on address: %s", input, address)

	command := []byte{DI_SETSVPERCENT}
	log.Printf("Command string: %s", hex.EncodeToString(command))

	firstDigit := strings.Split(address, ".")[1]
	nodeString := firstDigit[len(firstDigit)-1:] + strings.Split(address, ".")[2]

	log.Printf("Detected IP Address: %s", address)
	log.Printf("Detected HiQnet address: %s (decimal), %X (hex)", nodeString, (nodeString))

	nodeBytes, err := hex.DecodeString(nodeString)
	if err != nil {
		errorMessage := "Could not decode node string: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	for {
		if len(nodeBytes) > LEN_NODE {
			return []byte{}, errors.New("detected HiQnet node longer than 2 bytes")
		} else if len(nodeBytes) == LEN_NODE {
			nodeBytes = append([]byte{0x00}, nodeBytes...)
		}
	}

	command = append(command, nodeBytes...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, VIRTUAL_DEVICE)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	gainBlock, err := hex.DecodeString(input)
	if err != nil {
		errorMessage := "Could not decode input string: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	command = append(command, gainBlock...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, stateVariables["gain"]...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	log.Printf("Calculating parameter for volume %s", volume)
	toSend, _ := strconv.Atoi(volume)
	if toSend > 100 || toSend < 0 {
		return []byte{}, errors.New("Invalid volume request")
	}

	toSend *= 65536
	log.Printf("toSend: %v", toSend)

	hexValue := make([]byte, 4)
	binary.BigEndian.PutUint32(hexValue, uint32(toSend))

	command = append(command, hexValue...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	checksum := GetChecksumByte(command)

	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command, _ = MakeSubstitutions(command, ENCODE)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	STX := []byte{byte(ENCODE["STX"])}
	command = append(STX, command...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	ETX := []byte{byte(ENCODE["ETX"])}
	command = append(command, ETX...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	return command, nil
}
