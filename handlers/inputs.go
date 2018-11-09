package handlers

import (
	"net/http"
	"strconv"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/common/v2/auth"
	"github.com/byuoitav/london-audio-microservice/londondi"
	"github.com/labstack/echo"
)

const PORT = "1023"

func Mute(context echo.Context) error {
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "write-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

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

	command, err = londondi.Wrap(command)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	err = londondi.HandleRawCommandBytes(command, address+":"+PORT)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, status.Mute{true})
}

func UnMute(context echo.Context) error {
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "write-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

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

	command, err = londondi.Wrap(command)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	err = londondi.HandleRawCommandBytes(command, address+":"+PORT)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, status.Mute{false})
}

func SetVolume(context echo.Context) error {
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "write-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

	input := context.Param("input")
	address := context.Param("address")
	volume := context.Param("level")

	volumeInt, err := strconv.Atoi(volume)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	command, err := londondi.BuildRawVolumeCommand(input, address, volume)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	command, err = londondi.MakeSubstitutions(command, londondi.ENCODE)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	command, err = londondi.Wrap(command)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	err = londondi.HandleRawCommandBytes(command, address+":"+PORT)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, status.Volume{volumeInt})
}
