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

const NODE = 0x0001
const VIRTUAL_DEVICE = 0x03

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
