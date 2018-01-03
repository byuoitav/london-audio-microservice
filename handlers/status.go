package handlers

import (
	"fmt"
	"net/http"

	di "github.com/byuoitav/london-audio-microservice/londondi"
	"github.com/labstack/echo"
)

func GetMute(context echo.Context) error {

	status, err := di.GetMute(context.Param("address"), context.Param("input"))
	if err != nil {
		msg := fmt.Sprintf("Error: %s", err.Error())
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, status)

}

func GetVolume(context echo.Context) error {

	status, err := di.GetVolume(context.Param("address"), context.Param("input"))
	if err != nil {
		msg := fmt.Sprintf("Error: %s", err.Error())
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, status)
}
