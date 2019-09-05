package londondi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/byuoitav/common/pooled"
	"github.com/byuoitav/common/status"
	"github.com/fatih/color"
)

//GetVolume .
func GetVolume(address, input string) (status.Volume, error) {

	subscribe, err := BuildCommand(address, input, "volume", []byte{}, SubscribePercent)
	if err != nil {
		msg := fmt.Sprintf("unable to build subscribe command %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return status.Volume{}, errors.New(msg)
	}

	unsubscribe, err := BuildCommand(address, input, "volume", []byte{}, UnsubscribePercent)
	if err != nil {
		msg := fmt.Sprintf("unable to build unsubscribe command %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return status.Volume{}, errors.New(msg)
	}
	var response []byte
	work := func(conn pooled.Conn) error {
		response, err = GetStatus(subscribe, unsubscribe, address+":"+PORT, conn)
		if err != nil {
			msg := fmt.Sprintf("Could not execute commands: %s", err.Error())
			log.Printf("%s", color.HiRedString("[error] %s", msg))
			return errors.New(msg)
		}
		return nil
	}

	err = pool.Do(address, work)
	if err != nil {
		return status.Volume{}, err
	}

	response, err = Unwrap(response)
	if err != nil {
		errorMessage := "Could not unwrap message: " + err.Error()
		log.Printf(errorMessage)
		return status.Volume{}, errors.New(errorMessage)
	}

	response, err = MakeSubstitutions(response, DECODE)
	if err != nil {
		errorMessage := "Could not substitute reserved bytes: " + err.Error()
		log.Printf(errorMessage)
		return status.Volume{}, errors.New(errorMessage)
	}

	response, err = Validate(response)
	if err != nil {
		errorMessage := "Invalid message: " + err.Error()
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

//GetMute .
func GetMute(address, input string) (status.Mute, error) {

	log.Printf("%s", color.HiMagentaString("[status] getting mute status of channel %X from device at address %s", input, address))

	subscribe, err := BuildCommand(address, input, "mute", []byte{}, Subscribe)
	if err != nil {
		msg := fmt.Sprintf("unable to build subscribe command %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return status.Mute{}, errors.New(msg)
	}

	unsubscribe, err := BuildCommand(address, input, "mute", []byte{}, Unsubscribe)
	if err != nil {
		msg := fmt.Sprintf("unable to build unsubscribe command %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return status.Mute{}, errors.New(msg)
	}
	var response []byte
	work := func(conn pooled.Conn) error {
		response, err = GetStatus(subscribe, unsubscribe, address+":"+PORT, conn)
		if err != nil {
			errorMessage := "could not execute commands: " + err.Error()
			log.Printf(errorMessage)
			return errors.New(errorMessage)
		}
		return nil
	}

	err = pool.Do(address, work)
	if err != nil {
		return status.Mute{}, err
	}

	response, err = Unwrap(response)
	if err != nil {
		errorMessage := "Could not unwrap message: " + err.Error()
		log.Printf(errorMessage)
		return status.Mute{}, errors.New(errorMessage)
	}

	response, err = MakeSubstitutions(response, DECODE)
	if err != nil {
		errorMessage := "Could not substitute reserved bytes: " + err.Error()
		log.Printf(errorMessage)
		return status.Mute{}, errors.New(errorMessage)
	}

	response, err = Validate(response)
	if err != nil {
		errorMessage := "Invalid message: " + err.Error()
		log.Printf(errorMessage)
		return status.Mute{}, errors.New(errorMessage)
	}

	state, err := ParseMuteStatus(response)
	if err != nil {
		errorMessage := "Could not parse response: " + err.Error()
		log.Printf(errorMessage)
		return status.Mute{}, errors.New(errorMessage)
	}

	log.Printf("%s", color.HiMagentaString("[status] successfully retrieved status"))

	return state, nil

}

//BuildCommand .
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

//BuildRawCommand .
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

//GetStatus .
func GetStatus(subscribe, unsubscribe []byte, address string, pconn pooled.Conn) ([]byte, error) {

	log.Printf("[status] handling status command: %s...", color.HiMagentaString("%X", subscribe))

	log.Printf("[status] writing status command...")

	_, err := pconn.Write(subscribe)
	switch {
	case err != nil:
		return nil, fmt.Errorf("unable to subscribe: %s", err)
	}

	log.Printf("[status] reading status response...")
	reader := bufio.NewReader(pconn)
	response, err := reader.ReadBytes(ETX)
	if err != nil {
		msg := fmt.Sprintf("device not responding: %s", err.Error())
		log.Printf("%s", color.HiRedString("[error] %s", msg))
		return []byte{}, errors.New(msg)
	}

	log.Printf("[status] response: %s", color.HiBlueString("%x", response))
	log.Printf("[status] sending unsubscribe command: %s...", color.HiBlueString("%x", unsubscribe))

	_, err = pconn.Write(unsubscribe)
	switch {
	case err != nil:
		return nil, fmt.Errorf("unable to unsubscribe: %s", err)
	}

	return response, nil
}

//@pre: checksum byte removed
func ParseVolumeStatus(message []byte) (status.Volume, error) {

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

	return status.Volume{
		Volume: int(trueValue),
	}, nil
}

//@pre: checksum byte removed
func ParseMuteStatus(message []byte) (status.Mute, error) {

	log.Printf("Parsing mute status message: %X", message)

	//mute status determined with last byte
	data := message[len(message)-1:]
	log.Printf("data: %X", data)
	if bytes.EqualFold(data, []byte{0}) {
		return status.Mute{
			Muted: false,
		}, nil
	} else if bytes.EqualFold(data, []byte{1}) {
		return status.Mute{
			Muted: true,
		}, nil
	} else { //bad data
		return status.Mute{}, errors.New("Bad data in status message")
	}
}
