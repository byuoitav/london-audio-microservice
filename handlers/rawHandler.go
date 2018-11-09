package handlers

import (
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
	"github.com/byuoitav/london-audio-microservice/londondi"
	"github.com/jessemillar/jsonresp"
	"github.com/labstack/echo"
)

func Raw(context echo.Context) error {
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "write-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

	command := londondi.RawDICommand{}

	bindError := context.Bind(&command)
	if bindError != nil {
		jsonresp.New(context.Response(), http.StatusBadRequest, "Could not read command body: "+bindError.Error())
		return nil
	}

	commandError := londondi.HandleRawCommandString(command)
	if commandError != nil {
		jsonresp.New(context.Response(), http.StatusBadRequest, commandError.Error())
		return nil
	}

	return commandError
}

func RawInfo(context echo.Context) error {
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "write-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

	jsonresp.New(context.Response(), http.StatusBadRequest, "Send a POST request to the /raw endpoint with a body including Command string")
	return nil
}
