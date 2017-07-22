package londondi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
)

func HandleRawCommandString(rawCommand RawDICommand) error {

	log.Printf("Handling raw command: %s...", rawCommand.Command)

	connection, connectError := net.Dial("tcp", rawCommand.Address+":"+
		rawCommand.Port)
	if connectError != nil {
		log.Printf(connectError.Error())
		return connectError
	}

	log.Printf("Converting to command to hex value...")
	hexCommand, hexError := hex.DecodeString(rawCommand.Command)
	if hexError != nil {
		fmt.Println(hexError.Error())
		return hexError
	}

	log.Printf("hexCommand: %x", hexCommand)

	_, writeError := connection.Write(hexCommand)
	if writeError != nil {
		log.Printf(writeError.Error())
	}

	connection.Close()
	return connectError
}

func HandleRawCommandBytes(command []byte, address string) error {

	log.Printf("Handling raw command: %x...", command)

	connection, err := net.Dial("tcp", address)
	if err != nil {
		errorMessage := "Could not connect to device: " + err.Error()
		log.Printf(errorMessage)
		return errors.New(errorMessage)
	}

	_, err = connection.Write(command)
	if err != nil {
		errorMessage := "Could not write to device: " + err.Error()
		log.Printf(errorMessage)
		return errors.New(errorMessage)
	}

	connection.Close()
	return nil
}
