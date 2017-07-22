package handlers

import (
	"net/http"

	"github.com/byuoitav/london-audio-microservice/londondi"
	"github.com/labstack/echo"
)

const PORT = "1023"

func Mute(context echo.Context) error {

	input := context.Param("input")
	address := context.Param("address")

	command, err := londondi.BuildRawMuteCommand(input, address, "true")
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	command, err = londondi.MakeSubstitutions(command, londondi.ENCODE)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	command = londondi.Wrap(command)

	err = londondi.HandleRawCommandBytes(command, address+":"+PORT)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, "Success")
}

func UnMute(context echo.Context) error {

	input := context.Param("input")
	address := context.Param("address")

	command, err := londondi.BuildRawMuteCommand(input, address, "false")
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	command, err = londondi.MakeSubstitutions(command, londondi.ENCODE)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	command = londondi.Wrap(command)

	err = londondi.HandleRawCommandBytes(command, address+":"+PORT)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, "Success")
}

func SetVolume(context echo.Context) error {

	input := context.Param("input")
	address := context.Param("address")
	volume := context.Param("level")

	command, err := londondi.BuildRawVolumeCommand(input, address, volume)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	command, err = londondi.MakeSubstitutions(command, londondi.ENCODE)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	command = londondi.Wrap(command)

	err = londondi.HandleRawCommandBytes(command, address+":"+PORT)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, "Success")
}
