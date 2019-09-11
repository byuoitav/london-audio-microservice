package handlers

import (
	"net/http"
	"strconv"

	"github.com/byuoitav/common/status"
	"github.com/byuoitav/london-audio-microservice/londondi"
	"github.com/labstack/echo"
)

//Mute handler
func Mute(context echo.Context) error {
	input := context.Param("input")
	address := context.Param("address")

	err := londondi.Mute(input, address)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, status.Mute{true})
}

//UnMute handler
func UnMute(context echo.Context) error {
	input := context.Param("input")
	address := context.Param("address")

	err := londondi.UnMute(input, address)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, status.Mute{false})
}

//SetVolume handler
func SetVolume(context echo.Context) error {
	input := context.Param("input")
	address := context.Param("address")
	volume := context.Param("level")
	volumeInt, err := strconv.Atoi(volume)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	err = londondi.SetVolume(input, address, volume)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, status.Volume{volumeInt})
}
