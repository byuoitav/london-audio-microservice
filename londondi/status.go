package londondi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	se "github.com/byuoitav/av-api/statusevaluators"
	"github.com/byuoitav/london-audio-microservice/connect"
	"github.com/fatih/color"
)

func GetVolume(address, input string) (se.Volume, error) {

	subscribe, err := BuildCommand(address, input, "volume", []byte{}, SubscribePercent)
	if err != nil {
		msg := fmt.Sprintf("unable to build subscribe command %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return se.Volume{}, errors.New(msg)
	}

	unsubscribe, err := BuildCommand(address, input, "volume", []byte{}, UnsubscribePercent)
	if err != nil {
		msg := fmt.Sprintf("unable to build unsubscribe command %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return se.Volume{}, errors.New(msg)
	}

	response, err := GetStatus(subscribe, unsubscribe, address+":"+PORT)
	if err != nil {
		msg := fmt.Sprintf("Could not execute commands: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return se.Volume{}, errors.New(msg)
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

	log.Printf("%s", color.HiMagentaString("[status] getting mute status of channel %X from device at address %s", input, address))

	subscribe, err := BuildCommand(address, input, "mute", []byte{}, Subscribe)
	if err != nil {
		msg := fmt.Sprintf("unable to build subscribe command %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return se.MuteStatus{}, errors.New(msg)
	}

	unsubscribe, err := BuildCommand(address, input, "mute", []byte{}, Unsubscribe)
	if err != nil {
		msg := fmt.Sprintf("unable to build unsubscribe command %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return se.MuteStatus{}, errors.New(msg)
	}

	response, err := GetStatus(subscribe, unsubscribe, address+":"+PORT)
	if err != nil {
		errorMessage := "could not execute commands: " + err.Error()
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

	log.Printf("%s", color.HiMagentaString("[status] successfully retrieved status"))

	return state, nil

}

func BuildCommand(address, input, status string, data []byte, method Method) ([]byte, error) {

	log.Printf("[command] building command...")

	command, err := BuildRawCommand(address, input, status, data, method)
	if err != nil {
		msg := fmt.Sprintf("could not build subscribe command: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	command, err = MakeSubstitutions(command, ENCODE)
	if err != nil {
		msg := fmt.Sprintf("Could not substitute reserved bytes: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	command, err = Wrap(command)
	if err != nil {
		msg := fmt.Sprintf("Could not wrap message: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	return command, nil
}

func BuildRawCommand(address, input, state string, data []byte, method Method) ([]byte, error) {

	log.Printf("Building subscription message for %s on input %s at address %s", state, input, address)

	var base byte
	switch method {
	case Set:
		base = DI_SETSV
	case SetPercent:
		base = DI_SETSVPERCENT
	case Subscribe:
		base = DI_SUBSCRIBESV
	case Unsubscribe:
		base = DI_UNSUBSCRIBESV
	case SubscribePercent:
		base = DI_SUBSCRIBESVPERCENT
	case UnsubscribePercent:
		base = DI_UNSUBSCRIBESVPERCENT
	}

	command, err := GetCommandAddress(base, address)
	if err != nil {
		msg := fmt.Sprintf("could not address command: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	gainBlock, err := hex.DecodeString(input)
	if err != nil {
		msg := fmt.Sprintf("Could not decode input string: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	command = append(command, gainBlock...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	command = append(command, stateVariables[state]...)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	if method == Set || method == SetPercent {
		command = append(command, data...)
	} else if method == Subscribe || method == SubscribePercent {
		command = append(command, RATE...)
	} else { // it's an unsubscribe command and the rate is zero. I hope this works
		zero := make([]byte, len(RATE))
		command = append(command, zero...)
	}

	log.Printf("Command string: %s", hex.EncodeToString(command))

	checksum := GetChecksumByte(command)
	command = append(command, checksum)
	log.Printf("Command string: %s", hex.EncodeToString(command))

	return command, nil
}

func GetStatus(subscribe, unsubscribe []byte, address string) ([]byte, error) {

	log.Printf("[status] handling status command: %s...", color.HiMagentaString("%X", subscribe))

	connection, err := connect.GetConnection(address)
	if err != nil {
		msg := fmt.Sprintf("problem getting connection to device: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	log.Printf("[status] writing status command...")

	_, err = connection.Write(subscribe)
	if err != nil {
		msg := fmt.Sprintf("could not send subscribe message to device: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	log.Printf("[status] reading status response...")
	reader := bufio.NewReader(connection)
	response, err := reader.ReadBytes(ETX)

	if err != nil {
		msg := fmt.Sprintf("device not responding: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	log.Printf("[status] response: %s", color.HiBlueString("%x", response))
	log.Printf("[status] sending unsubscribe command: %s...", color.HiBlueString("%x", unsubscribe))

	_, err = connection.Write(unsubscribe)
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
