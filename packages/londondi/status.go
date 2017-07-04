package londondi

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"log"
	"net"

	"github.com/byuoitav/av-api/status"
)

func GetVolume(address, input string) (status.Volume, error) {

	command, err := BuildSubscribeCommand(address, input, "volume", DI_SUBSCRIBESVPERCENT)
	if err != nil {
		errorMessage := "Could not build subscribe command: " + err.Error()
		log.Printf(errorMessage)
		return status.Volume{}, errors.New(errorMessage)
	}

	response, err := HandleStatusCommand(command)
	if err != nil {
		errorMessage := "Could not execute commands: " + err.Error()
		log.Printf(errorMessage)
		return status.Volume{}, errors.New(errorMessage)
	}

	state, err := ParseVolumeStatus(response)
	if err != nil {
		errorMessage := "Could not parse response: " + err.Error()
		log.Printf(errorMessage)
		return status.Volume{}, errors.New(errorMessage)
	}

	return state, nil

}

func GetMute(address, input string) (status.MuteStatus, error) {

	command, err := BuildSubscribeCommand(address, input, "mute", DI_SUBSCRIBESV)
	if err != nil {
		errorMessage := "Could not build subscribe command: " + err.Error()
		log.Printf(errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	response, err := HandleStatusCommand(command)
	if err != nil {
		errorMessage := "Could not execute commands: " + err.Error()
		log.Printf(errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	state, err := ParseMuteStatus(response)
	if err != nil {
		errorMessage := "Could not parse response: " + err.Error()
		log.Printf(errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	return state, nil

}

func BuildSubscribeCommand(address, input, state string, messageType int32) (RawDICommand, error) {

	log.Printf("Building raw command to subsribe to %s of input %s on address %s", state, input, address)

	command := []byte{byte(messageType)}

	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, NODE...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	object, _ := hex.DecodeString(gainBlocks[input])
	command = append(command, object...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	stateVariable, _ := hex.DecodeString(stateVariables[state])
	command = append(command, stateVariable...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, RATE...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	checksum := GetChecksumByte(command)
	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command, _ = MakeSubstitutions(command, reserved)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	STX := []byte{byte(reserved["STX"])}
	command = append(STX, command...)
	ETX := byte(reserved["ETX"])
	command = append(command, ETX)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	return RawDICommand{
		Address: address,
		Port:    PORT,
		Command: hex.EncodeToString(command),
	}, nil
}

func HandleStatusCommand(subscribe RawDICommand) ([]byte, error) {

	log.Printf("Handling status command...")

	connection, err := net.Dial("tcp", subscribe.Address+":"+subscribe.Port)
	if err != nil {
		errorMessage := "Could not connect to device: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	defer connection.Close()

	log.Printf("Converting command to hex value...")
	command, err := hex.DecodeString(subscribe.Command)
	if err != nil {
		errorMessage := "Could not convert command to hex value: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	_, err = connection.Write(command)
	if err != nil {
		errorMessage := "Could not send message to device: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	reader := bufio.NewReader(connection)

	response, err := reader.ReadBytes(ETX)
	if err != nil {
		errorMessage := "Could not find ETX: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	log.Printf("Message: %x", response)

	return response, nil

}

func ParseVolumeStatus(message []byte) (status.Volume, error) {

	return status.Volume{}, nil
}

func ParseMuteStatus(message []byte) (status.MuteStatus, error) {

	log.Printf("Parsing message: %X", message)
	decoded, err := RemoveEscapeCharacters(message)
	if err != nil {
		errorMessage := "Could not remove escape characters: " + err.Error()
		log.Printf(errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	log.Printf("Without escape characters: %X", decoded)

	buffer := bytes.NewBuffer(decoded)
	stx, err := buffer.ReadByte()
	if err != nil {
		errorMessage := "Could not read byte: " + err.Error()
		log.Printf(errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	if stx != STX {
		errorMessage := "Status does not start with STX byte."
		log.Printf("Error: %s", errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	messageType, err := buffer.ReadByte()
	if err != nil {
		errorMessage := "Could not read message type byte"
		log.Printf("Error: %s", errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	if messageType != DI_SETSV {
		errorMessage := "Status does not start with STX byte."
		log.Printf("Error: %s", errorMessage)
		return status.MuteStatus{}, errors.New(errorMessage)
	}

	return status.MuteStatus{}, nil
}

func RemoveEscapeCharacters(message []byte) ([]byte, error) {

	log.Printf("Removing escape characters...")

	return []byte{}, nil
}
