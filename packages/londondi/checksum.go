package londondi

import (
	"bytes"
	"encoding/hex"
	"log"
)

func GetChecksumByte(command []byte) byte {

	log.Printf("Generating checksum byte...")
	log.Printf("command %v", command)
	log.Printf("message length: %v", len(command))

	checksum := command[0] ^ command[1]

	for i := 2; i < len(command); i++ {
		checksum = checksum ^ command[i]
	}

	return checksum
}

func MakeSubstitutions(command []byte) ([]byte, error) {

	log.Printf("Making substitutions...")

	newCommand, err := FindAndReplace(command, 0x1b)
	if err != nil {
		return []byte{}, err
	}

	command = newCommand

	for key, value := range reserved {

		if key == "Escape" {
			continue
		}

		newCommand, err := FindAndReplace(command, value)
		if err != nil {
			return []byte{}, err
		}

		command = newCommand
	}

	return command, nil
}

func FindAndReplace(command []byte, reserved int) ([]byte, error) {

	fragments := bytes.Split(command, []byte{byte(reserved)})
	if len(fragments) > 1 {

		newBytes, err := hex.DecodeString(substitutions[reserved])
		if err != nil {
			return []byte{}, err
		}

		return bytes.Join(fragments, newBytes), nil
	}

	return command, nil
}
