package londondi

import (
	"encoding/hex"
	"log"
)

/* to set a state variable, <DI_SETSV>
<node>
<virtual_device>
<object>
<state_variable>
<data>
*/

func BuildRawMuteCommand(input, address string) (RawDICommand, error) {

	log.Printf("Building raw mute command for input: %s on address: %s", input, address)

	command := []byte{byte(DI_SETSV)}
	log.Printf("Command string: %s", hex.EncodeToString(command))

	stupid := append([]byte{byte(0x00)}, byte(NODE))
	command = append(command, stupid...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, []byte{byte(VIRTUAL_DEVICE)}...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	object, _ := hex.DecodeString(gainBlocks[input])
	command = append(command, object...)

	stateVariable, _ := hex.DecodeString(stateVariables["mute"])
	command = append(command, stateVariable...)

	data, _ := hex.DecodeString(muteStates["true"])
	command = append(command, data...)

	checksum := GetChecksumByte(command)

	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command, _ = MakeSubstitutions(command)
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

func BuildRawUnMuteCommand(input string, address string) (RawDICommand, error) {

	log.Printf("Building raw unmute command for input: %s on address: %s", input, address)

	command := []byte{byte(DI_SETSV)}
	log.Printf("Command string: %s", hex.EncodeToString(command))

	stupid := append([]byte{byte(0x00)}, byte(NODE))
	command = append(command, stupid...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, []byte{byte(VIRTUAL_DEVICE)}...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	object, _ := hex.DecodeString(gainBlocks[input])
	command = append(command, object...)

	stateVariable, _ := hex.DecodeString(stateVariables["mute"])
	command = append(command, stateVariable...)

	data, _ := hex.DecodeString(muteStates["false"])
	command = append(command, data...)

	checksum := GetChecksumByte(command)

	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command, _ = MakeSubstitutions(command)
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

	command := []byte{byte(DI_SETSV)}
	log.Printf("Command string: %s", hex.EncodeToString(command))

	stupid := append([]byte{byte(0x00)}, byte(NODE))
	command = append(command, stupid...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, []byte{byte(VIRTUAL_DEVICE)}...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	object, _ := hex.DecodeString(gainBlocks[input])
	command = append(command, object...)

	stateVariable, _ := hex.DecodeString(stateVariables["mute"])
	command = append(command, stateVariable...)

	data, _ := hex.DecodeString(muteStates["false"])
	command = append(command, data...)

	checksum := GetChecksumByte(command)

	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command, _ = MakeSubstitutions(command)
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
