package londondi

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

func BuildSubscribeCommand(address, input, state string, messageType byte) (RawDICommand, error) {

	log.Printf("Building raw command to subsribe to %s of input %s on address %s", state, input, address)

	command := []byte{byte(messageType)}

	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, NODE...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	gainBlock, err := hex.DecodeString(input)
	if err != nil {
		errorMessage := "Could not decode input string: " + err.Error()
		log.Printf(errorMessage)
		return RawDICommand{}, errors.New(errorMessage)
	}

	command = append(command, gainBlock...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, stateVariables[state]...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, RATE...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	checksum := GetChecksumByte(command)
	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command, _ = MakeSubstitutions(command, ENCODE)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	stx := []byte{STX}
	command = append(stx, command...)
	ETX := byte(ENCODE["ETX"])
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

	log.Printf("Response: %x", response)

	return response, nil

}

func ParseVolumeStatus(message []byte) (status.Volume, error) {

	message, err := ValidateMessage(message)
	if err != nil {
		return status.Volume{}, err
	}

	//get data - always 4 bytes
	data := message[len(message)-4:]
	log.Printf("data: %X", data)
	log.Printf("len(data): %v", len(data))

	//turn data into number between 0 and 100
	const SCALE_FACTOR = 65536
	var rawValue int32
	_ = binary.Read(bytes.NewReader(data), binary.BigEndian, &rawValue)
	log.Printf("rawValue %v", rawValue)

	trueValue := rawValue / SCALE_FACTOR

	trueValue++ //not sure why it comes up with the wrong number

	return status.Volume{
		Volume: int(trueValue),
	}, nil
}

func ParseMuteStatus(message []byte) (status.MuteStatus, error) {

	log.Printf("Parsing mute status message: %X", message)

	message, err := ValidateMessage(message)
	if err != nil {
		return status.MuteStatus{}, err
	}

	data := message[len(message)-1:]
	log.Printf("data: %X", data)
	if bytes.EqualFold(data, []byte{0}) {
		return status.MuteStatus{
			Muted: false,
		}, nil
	} else if bytes.EqualFold(data, []byte{1}) {
		return status.MuteStatus{
			Muted: true,
		}, nil
	} else { //bad data
		return status.MuteStatus{}, errors.New("Bad data in status message")
	}
}

//validates message and returns the message with STX, ETX, and checksum bytes removed
func ValidateMessage(message []byte) ([]byte, error) {

	log.Printf("Validating status message %X", message)

	//remove STX
	message = bytes.TrimPrefix(message, []byte{STX})
	log.Printf("message %X", message)

	//remove ETX
	message = bytes.TrimSuffix(message, []byte{ETX})
	log.Printf("message %X", message)

	//make substitutions
	message, err := MakeSubstitutions(message, DECODE)
	if err != nil {
		errorMessage := "Could not make substitutions" + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	log.Printf("message %X", message)
	//grrr...
	message = bytes.Replace(message, []byte{0x1b}, []byte{}, -1)

	log.Printf("message %X", message)

	//check checksum
	checksum := GetChecksumByte(message[:len(message)-1])
	if checksum != message[len(message)-1] {
		errorMessage := "Checksums do not match"
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	message = bytes.TrimSuffix(message, []byte{checksum})
	log.Printf("message %X", message)

	return message, nil
}
