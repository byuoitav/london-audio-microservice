package londondi

type RawDICommand struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Command string `json:"command"`
}

type RawDIResponse struct {
	Response []string `json:"response"`
}

//standardize file based on these values
var gainBlocks = map[string]string{

	"mic1":   "00010a",
	"mic2":   "00010b",
	"mic3":   "00010c",
	"mic4":   "00010d",
	"media1": "00010e",
	"media2": "00010f",
}

//standards in London documentation
var stateVariables = map[string]string{

	"gain":     "0000",
	"mute":     "0001",
	"polarity": "0002",
}

var muteStates = map[string]string{

	"true":  "00000001",
	"false": "00000000",
}

var PORT = "1023"

const DI_SETSV = 0x88
const DI_SETSVPERCENT = 0x8d
const DI_SUBSCRIBESV = 0x89
const DI_SUBSCRIBESVPERCENT = 0x8e
const DI_UNSUBSCRIBESV = 0x8a
const DI_UNSUBSCRIBESVPERCENT = 0x8f

//2 bytes for NODE, 1 byte for VIRTUAL_DEVICE should be the same for all cases!
var NODE = []byte{0x00, 0x01, 0x03}

var RATE = []byte{0x00, 0x00, 0x00, 0x32} //represents 50 ms, the shortest interval

var ACK = byte(0x06)
var ETX = byte(0x03)

var reserved = map[string]int{
	"STX":    0x02,
	"ETX":    0x03,
	"ACK":    0x06,
	"NAK":    0x15,
	"Escape": 0x1b,
}

var substitutions = map[int]string{

	0x02: "1b82",
	0x03: "1b83",
	0x06: "1b86",
	0x15: "1b95",
	0x1b: "1b9b",
}
