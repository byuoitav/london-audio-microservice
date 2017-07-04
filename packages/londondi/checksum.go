package londondi

import (
	"bytes"
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

func MakeSubstitutions(command []byte, toCheck map[string]int) ([]byte, error) {

	log.Printf("Making substitutions...")

	//always address escape byte first
	escapeInt := toCheck["escape"]
	escapeByte := byte(escapeInt)

	log.Printf("replacing %x with %x", escapeInt, substitutions[escapeInt])
	newCommand := bytes.Replace(command, []byte{escapeByte}, substitutions[escapeInt], -1)

	for key, value := range toCheck {

		if key == "escape" {
			continue
		}

		log.Printf("replacing %x with %x", value, substitutions[value])
		newCommand = bytes.Replace(newCommand, []byte{byte(value)}, substitutions[value], -1)
		log.Printf("changed command: %x", newCommand)

	}

	return newCommand, nil
}
