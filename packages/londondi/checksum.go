package londondi

import (
	"bytes"
	"log"
)

func GetChecksumByte(message []byte) byte {

	log.Printf("Generating checksum byte for message %x...", message)

	checksum := message[0] ^ message[1]

	for i := 2; i < len(message); i++ {
		checksum = checksum ^ message[i]
	}

	log.Printf("checksum: %x", checksum)
	return checksum
}

func MakeSubstitutions(command []byte, toCheck map[string]int) ([]byte, error) {

	log.Printf("Making substitutions...")

	//always address escape byte first
	escapeInt := toCheck["escape"]
	escapeByte := byte(escapeInt)

	log.Printf("replacing %x with %x", escapeInt, substitutions[escapeInt])
	newCommand := bytes.Replace(command, []byte{escapeByte}, substitutions[escapeInt], -1)
	log.Printf("changed command: %x", newCommand)

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
