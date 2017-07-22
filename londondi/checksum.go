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
	var toReplace []byte
	toReplace = append(toReplace, byte(escapeInt))

	if len(substitutions[escapeInt]) == 1 {
		//get the second bit
		newEscapeInt := escapeInt >> 8
		toReplace = append([]byte{byte(newEscapeInt)}, toReplace...)
	}

	log.Printf("replacing %x with %x", toReplace, substitutions[escapeInt])
	newCommand := bytes.Replace(command, toReplace, substitutions[escapeInt], -1)
	log.Printf("changed command: %x", newCommand)

	for key, value := range toCheck {

		if key == "escape" {
			continue
		}

		var iHateYou []byte
		iHateYou = append(iHateYou, byte(value))

		if len(substitutions[value]) == 1 {
			//get the second bit
			newEscapeInt := value >> 8
			iHateYou = append([]byte{byte(newEscapeInt)}, iHateYou...)
		}
		log.Printf("replacing %x with %x", value, substitutions[value])
		newCommand = bytes.Replace(newCommand, iHateYou, substitutions[value], -1)
		log.Printf("changed command: %x", newCommand)

	}

	stx := []byte{STX}
	newCommand = append(stx, newCommand...)
	newCommand = append(newCommand, ETX)
	return newCommand, nil
}
