package londondi

import (
	"encoding/hex"
	"fmt"
	"net"
)

func HandleRawCommand(rawCommand RawDICommand) error {
	connection, connectError := net.Dial("tcp", rawCommand.Address+":"+
		rawCommand.Port)
	if connectError != nil {
		fmt.Println(connectError.Error())
		return connectError
	}

	hexCommand, hexError := hex.DecodeString(rawCommand.Command)
	if hexError != nil {
		fmt.Println(hexError.Error())
		return hexError
	}

	_, writeError := connection.Write(hexCommand)
	if writeError != nil {
		fmt.Println(writeError.Error())
	}

	connection.Close()
	return connectError
}
