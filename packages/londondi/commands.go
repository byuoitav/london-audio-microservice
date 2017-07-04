package londondi

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"strconv"
)

/* to set a state variable, <DI_SETSV>
<node>
<virtual_device>
<object>
<state_variable>
<data>
*/

func BuildRawMuteCommand(input, address, status string) (RawDICommand, error) {

	log.Printf("Building raw mute command for input: %s on address: %s", input, address)

	command := []byte{byte(DI_SETSV)}
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, NODE...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	object, _ := hex.DecodeString(gainBlocks[input])
	command = append(command, object...)

	stateVariable, _ := hex.DecodeString(stateVariables["mute"])
	command = append(command, stateVariable...)

	data, _ := hex.DecodeString(muteStates[status])
	command = append(command, data...)

	checksum := GetChecksumByte(command)

	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command, _ = MakeSubstitutions(command, reserved)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	STX := []byte{byte(reserved["STX"])}
	command = append(STX, command...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	ETX := []byte{byte(reserved["ETX"])}
	command = append(command, ETX...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	//since we're building a mute command, we a mute state variable for the specific port
	return RawDICommand{
		Address: address,
		Port:    PORT,
		Command: hex.EncodeToString(command),
	}, nil

}

func BuildRawVolumeCommand(input string, address string, volume string) (RawDICommand, error) {

	log.Printf("Building raw volume command for input: %s on address: %s", input, address)

	command := []byte{byte(DI_SETSVPERCENT)}
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, NODE...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	object, _ := hex.DecodeString(gainBlocks[input])
	command = append(command, object...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	state, _ := hex.DecodeString(stateVariables["gain"])
	command = append(command, state...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	log.Printf("Calculating parameter for volume %s", volume)
	toSend, _ := strconv.Atoi(volume)
	if toSend > 100 || toSend < 0 {
		return RawDICommand{}, errors.New("Invalid volume request")
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

	command, _ = MakeSubstitutions(command, reserved)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	STX := []byte{byte(reserved["STX"])}
	command = append(STX, command...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	ETX := []byte{byte(reserved["ETX"])}
	command = append(command, ETX...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	//since we're building a mute command, we a mute state variable for the specific port
	return RawDICommand{
		Address: address,
		Port:    PORT,
		Command: hex.EncodeToString(command),
	}, nil

}
