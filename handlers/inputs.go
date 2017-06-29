package handlers

import (
	"net/http"

	"github.com/byuoitav/london-audio-microservice/packages/londondi"
	"github.com/labstack/echo"
)

func Mute(context echo.Context) error {

	commands, err := GenerateMuteCommands(context.Param("input"), context.Param("address"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	err = londondi.HandleRawCommands(commands)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, "success")
}

func UnMute(context echo.Context) error {

	input := context.Param("input")
	address := context.Param("address")

	command, err := londondi.BuildRawUnMuteCommand(input, address)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	err = londondi.HandleRawCommand(command)
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

	err = londondi.HandleRawCommand(command)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, "Success")
}
