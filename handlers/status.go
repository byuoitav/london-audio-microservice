package handlers

import (
	"net/http"

	di "github.com/byuoitav/london-audio-microservice/londondi"
	"github.com/labstack/echo"
)

func GetMute(context echo.Context) error {

	status, err := di.GetMute(context.Param("address")+":"+PORT, context.Param("input"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, status)

}

func GetVolume(context echo.Context) error {

	status, err := di.GetVolume(context.Param("address")+":"+PORT, context.Param("input"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, status)
}
