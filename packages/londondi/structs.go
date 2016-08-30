package londondi

type RawDICommand struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Command string `json:"command"`
}

type RawDIResponse struct {
	Response []string `json:"response"`
}
