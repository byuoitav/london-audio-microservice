package londondi

type RawDICommand struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Command string `json:"command"`
}

type RawDIResponse struct {
	Response []string `json:"response"`
}

//london DI commands in semi-readable format

var tokens = map[string]string{

	"STX": "02",
	"ETX": "03",
}

var commands = map[string]string{

	"DI_SETSV":        "88",
	"DI_SETSVPERCENT": "8D",
}

var constants = map[string]string{

	"node":          "0001",
	"virtualDevice": "03",
}

var cards = map[string]string{

	"mic":   "000001",
	"media": "000002",
}

var gains = map[string]string{

	"mic1":   "4",
	"mic2":   "a",
	"mic3":   "10",
	"mic4":   "16",
	"media1": "4",
	"media2": "10",
	"media3": "a",
	"media4": "16",
}

var mutes = map[string]string{

	"mic1":   "7d0",
	"mic2":   "7d1",
	"mic3":   "7d2",
	"mic4":   "7d3",
	"media1": "7d0",
	"media2": "7d2",
	"media3": "7d1",
	"media4": "7d3",
}
