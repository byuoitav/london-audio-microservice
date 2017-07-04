package londondi

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

func HandleRawCommand(rawCommand RawDICommand) error {

	log.Printf("Handling raw command: %s...", rawCommand.Command)

	connection, connectError := net.Dial("tcp", rawCommand.Address+":"+
		rawCommand.Port)
	if connectError != nil {
		fmt.Println(connectError.Error())
		return connectError
	}

	log.Printf("Converting to command to hex value...")
	hexCommand, hexError := hex.DecodeString(rawCommand.Command)
	if hexError != nil {
		fmt.Println(hexError.Error())
		return hexError
	}

	log.Printf("hexCommand: %v", hexCommand)

	_, writeError := connection.Write(hexCommand)
	if writeError != nil {
		fmt.Println(writeError.Error())
	}

	connection.Close()
	return connectError
}
