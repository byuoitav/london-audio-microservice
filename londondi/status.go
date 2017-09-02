package londondi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"

	se "github.com/byuoitav/av-api/statusevaluators"
	"github.com/byuoitav/london-audio-microservice/connect"
	"github.com/fatih/color"
)

const TIMEOUT = 10

func GetVolume(address, input string) (se.Volume, error) {

	command, err := BuildSubscribeCommand(address, input, "volume", DI_SUBSCRIBESVPERCENT)
	if err != nil {
		errorMessage := "Could not build subscribe command: " + err.Error()
		log.Printf(errorMessage)
		return se.Volume{}, errors.New(errorMessage)
	}

	command, err = MakeSubstitutions(command, ENCODE)
	if err != nil {
		errorMessage := "Could not substitute escape bytes: " + err.Error()
		log.Printf(errorMessage)
		return se.Volume{}, errors.New(errorMessage)
	}

	command, err = Wrap(command)
	if err != nil {
		errorMessage := "Could not wrap command: " + err.Error()
		log.Printf(errorMessage)
		return se.Volume{}, errors.New(errorMessage)
	}

	response, err := HandleStatusCommand(command, address+":"+PORT)
	if err != nil {
		errorMessage := "Could not execute commands: " + err.Error()
		log.Printf(errorMessage)
		return se.Volume{}, errors.New(errorMessage)
	}

	response, err = Unwrap(response)
	if err != nil {
		errorMessage := "Could not unwrap message: " + err.Error()
		log.Printf(errorMessage)
		return se.Volume{}, errors.New(errorMessage)
	}

	response, err = MakeSubstitutions(response, DECODE)
	if err != nil {
		errorMessage := "Could not substitute reserved bytes: " + err.Error()
		log.Printf(errorMessage)
		return se.Volume{}, errors.New(errorMessage)
	}

	response, err = Validate(response)
	if err != nil {
		errorMessage := "Invalid message: " + err.Error()
		log.Printf(errorMessage)
		return se.Volume{}, errors.New(errorMessage)
	}

	state, err := ParseVolumeStatus(response)
	if err != nil {
		errorMessage := "Could not parse response: " + err.Error()
		log.Printf(errorMessage)
		return se.Volume{}, errors.New(errorMessage)
	}

	return state, nil

}

func GetMute(address, input string) (se.MuteStatus, error) {

	command, err := BuildSubscribeCommand(address, input, "mute", DI_SUBSCRIBESV)
	if err != nil {
		errorMessage := "Could not build subscribe command: " + err.Error()
		log.Printf(errorMessage)
		return se.MuteStatus{}, errors.New(errorMessage)
	}

	command, err = MakeSubstitutions(command, ENCODE)
	if err != nil {
		errorMessage := "Could not substitute reserved bytes: " + err.Error()
		log.Printf(errorMessage)
		return se.MuteStatus{}, errors.New(errorMessage)
	}

	command, err = Wrap(command)
	if err != nil {
		errorMessage := "Could not wrap command: " + err.Error()
		log.Printf(errorMessage)
		return se.MuteStatus{}, errors.New(errorMessage)
	}

	response, err := HandleStatusCommand(command, address+":"+PORT)
	if err != nil {
		errorMessage := "Could not execute commands: " + err.Error()
		log.Printf(errorMessage)
		return se.MuteStatus{}, errors.New(errorMessage)
	}

	response, err = Unwrap(response)
	if err != nil {
		errorMessage := "Could not unwrap message: " + err.Error()
		log.Printf(errorMessage)
		return se.MuteStatus{}, errors.New(errorMessage)
	}

	response, err = MakeSubstitutions(response, DECODE)
	if err != nil {
		errorMessage := "Could not substitute reserved bytes: " + err.Error()
		log.Printf(errorMessage)
		return se.MuteStatus{}, errors.New(errorMessage)
	}

	response, err = Validate(response)
	if err != nil {
		errorMessage := "Invalid message: " + err.Error()
		log.Printf(errorMessage)
		return se.MuteStatus{}, errors.New(errorMessage)
	}

	state, err := ParseMuteStatus(response)
	if err != nil {
		errorMessage := "Could not parse response: " + err.Error()
		log.Printf(errorMessage)
		return se.MuteStatus{}, errors.New(errorMessage)
	}

	return state, nil

}

func BuildSubscribeCommand(address, input, state string, messageType byte) ([]byte, error) {

	log.Printf("Building raw command to subsribe to %s of input %s on address %s", state, input, address)

	command, err := GetCommandAddress(messageType, address)
	if err != nil {
		errorMessage := "Could not address command: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	gainBlock, err := hex.DecodeString(input)
	if err != nil {
		errorMessage := "Could not decode input string: " + err.Error()
		log.Printf(errorMessage)
		return []byte{}, errors.New(errorMessage)
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

	return command, nil
}

func HandleStatusCommand(subscribe, unsubscribe []byte, address string) ([]byte, error) {

	log.Printf("[status] handling status command: %s...", color.HiBlueString("%x", subscribe))

	connection, err := connection.GetConnection(address)
	if err != nil {
		msg := fmt.Sprintf("problem getting connection to device: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	_, err = connection.Write(subscribe)
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		_, err := connect.HandleTimeout(&connection, subscribe, connect.Write)
	}

	if err != nil {
		msg := fmt.Sprintf("could not send subscribe message to device: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	reader := bufio.NewReader(connection)
	response, err := reader.ReadBytes(ETX)
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		response, err = HandleTimeout(connection, []byte{ETX}, connect.Read)
	}

	if err != nil {
		msg := color.HiRedString("could not read status message to device: %s", err.Error())
		log.Printf("[error] %s", errorMessage)
		return []byte{}, errors.New(errorMessage)
	}

	log.Printf("[status] response: %s", color.HiBlueString("%x", response))
	log.Printf("[status] sending unsubscribe command: %s...", color.HiBlueString("%x", unsubscribe))

	_, err = connection.Write(unsubscribe)
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		_, err = connect.HandleTimeout(connection, unsubscribe, connect.Write)
	}

	if err != nil {
		msg := fmt.Sprintf("could not	not send unsubscribe message to device: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	return response, nil
}

//@pre: checksum byte removed
func ParseVolumeStatus(message []byte) (se.Volume, error) {

	log.Printf("Parsing mute status message: %X", message)

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

	return se.Volume{
		Volume: int(trueValue),
	}, nil
}

//@pre: checksum byte removed
func ParseMuteStatus(message []byte) (se.MuteStatus, error) {

	log.Printf("Parsing mute status message: %X", message)

	//mute status determined with last byte
	data := message[len(message)-1:]
	log.Printf("data: %X", data)
	if bytes.EqualFold(data, []byte{0}) {
		return se.MuteStatus{
			Muted: false,
		}, nil
	} else if bytes.EqualFold(data, []byte{1}) {
		return se.MuteStatus{
			Muted: true,
		}, nil
	} else { //bad data
		return se.MuteStatus{}, errors.New("Bad data in status message")
	}
}

func handleTimeout(conn *net.TCPConn, msg []byte, dir string)
