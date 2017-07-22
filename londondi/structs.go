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
var gainBlocks = map[string][]byte{

	"mic1":   {0x00, 0x01, 0x0a},
	"mic2":   {0x00, 0x01, 0x0b},
	"mic3":   {0x00, 0x01, 0x0c},
	"mic4":   {0x00, 0x01, 0x0d},
	"media1": {0x00, 0x01, 0x0e},
	"media2": {0x00, 0x01, 0x0f},
}

//standards in London documentation
var stateVariables = map[string][]byte{

	"gain":     {0x00, 0x00},
	"mute":     {0x00, 0x01},
	"polarity": {0x00, 0x02},
}

var muteStates = map[string][]byte{

	"true":  {0x00, 0x00, 0x00, 0x01},
	"false": {0x00, 0x00, 0x00, 0x00},
}

var test = map[bool][]byte{

	true:  {0x00, 0x00, 0x00, 0x01},
	false: {0x00, 0x00, 0x00, 0x00},
}

var PORT = "1023"
var DI_SETSV = byte(0x88)
var DI_SETSVPERCENT = byte(0x8d)
var DI_SUBSCRIBESV = byte(0x89)
var DI_SUBSCRIBESVPERCENT = byte(0x8e)
var DI_UNSUBSCRIBESV = byte(0x8a)
var DI_UNSUBSCRIBESVPERCENT = byte(0x8f)

//2 bytes for NODE, 1 byte for VIRTUAL_DEVICE should be the same for all cases!
var VIRTUAL_DEVICE = byte(0x03)

var RATE = []byte{0x00, 0x00, 0x00, 0x32} //represents 50 ms, the shortest interval

var ACK = byte(0x06)
var ETX = byte(0x03)
var STX = byte(0x02)

var ENCODE = map[string]int{
	"STX":    0x02,
	"ETX":    0x03,
	"ACK":    0x06,
	"NAK":    0x15,
	"escape": 0x1b,
}

var DECODE = map[string]int{
	"STX":    0x1b82,
	"ETX":    0x1b83,
	"ACK":    0x1b86,
	"NAK":    0x1b95,
	"escape": 0x1b9b,
}

var substitutions = map[int][]byte{

	0x02:   {0x1b, 0x82},
	0x03:   {0x1b, 0x83},
	0x06:   {0x1b, 0x86},
	0x15:   {0x1b, 0x95},
	0x1b:   {0x1b, 0x9b},
	0x1b82: {0x02},
	0x1b83: {0x03},
	0x1b86: {0x06},
	0x1b95: {0x15},
	0x1b9b: {0x1b},
}
