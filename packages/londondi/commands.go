package londondi

import (
	"errors"
	"strings"
)

/* to set a state variable, <DI_SETSV>
<node>
<virtual_device>
<object>
<state_variable>
<data>
*/

var PORT = "1023"

func BuildRawMuteCommand(input, address string) (RawDICommand, error) {

	commandString := tokens["STX"] + commands["DI_SETSV"] + constants["node"] + constants["virtualDevice"]

	//get object based on given input
	if strings.Contains(input, "mic") {

		commandString = commandString + cards["mic"]

	} else if strings.Contains(input, "media") {

		commandString = commandString + cards["media"]

	} else {
		return RawDICommand{}, errors.New("Invalid port")
	}

	commandString = commandString + mutes[input]

	//since we're building a mute command, we a mute state variable for the specific port
	return RawDICommand{
		Address: address,
		Port:    PORT,
		Command: commandString,
	}, nil

}

func BuildRawUnMuteCommand(input string, address string) (RawDICommand, error) {

	return RawDICommand{}, nil

}

func BuildRawVolumeCommand(input string, address string, volume string) (RawDICommand, error) {

	return RawDICommand{}, nil

}
