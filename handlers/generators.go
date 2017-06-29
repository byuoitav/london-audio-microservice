package handlers

import (
	"errors"
	"strconv"
	"strings"

	di "github.com/byuoitav/london-audio-microservice/packages/londondi"
)

func GenerateMuteCommands(input, address string) ([]di.RawDICommand, error) {

	if strings.Contains(input, "media") { //generate two commands for stereo inputs
		//gains and mutes mapped such that incrementing the number by two corresponds to the same physical input

		command1, err := di.BuildRawMuteCommand(input, address)
		if err != nil {
			return []di.RawDICommand{}, err
		}

		number, _ := strconv.Atoi(input[1:])
		input2 := input[:len(input)-1] + strconv.Itoa(number+2)

		command2, err := di.BuildRawMuteCommand(input2, address)
		if err != nil {
			return []di.RawDICommand{}, err
		}

		return []di.RawDICommand{command1, command2}, nil

	} else if strings.Contains(input, "mic") { //mic ports are mono

		command, err := di.BuildRawMuteCommand(input, address)
		if err != nil {
			return []di.RawDICommand{}, err
		}

		return []di.RawDICommand{command}, nil

	} else { // bad input

		return []di.RawDICommand{}, errors.New("Invalid port.")
	}

}
