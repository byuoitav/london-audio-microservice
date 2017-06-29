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

var commands = map[string]string{

	"DI_SETSV":        "88",
	"DI_SETSVPERCENT": "8D",
}

var tokens = map[string]string{

	"STX": "02",
	"ETX": "03",
}
