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
		jsonresp.New(context.Response(), http.StatusBadRequest, "Could not read command body: "+bindError.Error())
		return nil
	}

	commandError := londondi.HandleRawCommand(command)
	if commandError != nil {
		jsonresp.New(context.Response(), http.StatusBadRequest, commandError.Error())
		return nil
	}

	return commandError
}

func RawInfo(context echo.Context) error {
	jsonresp.New(context.Response(), http.StatusBadRequest, "Send a POST request to the /raw endpoint with a body including Command string")
	return nil
}
