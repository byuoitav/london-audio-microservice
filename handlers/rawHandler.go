package handlers

import (
	"net/http"

	"github.com/byuoitav/london-audio-microservice/packages/londondi"
	"github.com/jessemillar/jsonresp"
	"github.com/labstack/echo"
)

func Raw(context echo.Context) error {
	command := londondi.RawDICommand{}

	bindError := context.Bind(&command)
	if bindError != nil {
		return jsonresp.New(context, http.StatusBadRequest, "Could not read command body: "+bindError.Error())
	}

	commandError := londondi.HandleRawCommand(command)
	if commandError != nil {
		return jsonresp.New(context, http.StatusBadRequest, commandError.Error())
	}
	return commandError
}

func RawInfo(context echo.Context) error {
	return jsonresp.New(context, http.StatusBadRequest, "Send a POST request to the /raw endpoint with a body including Command string")
}
