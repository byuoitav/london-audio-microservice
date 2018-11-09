package handlers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
	di "github.com/byuoitav/london-audio-microservice/londondi"
	"github.com/labstack/echo"
)

func GetMute(context echo.Context) error {
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "read-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

	status, err := di.GetMute(context.Param("address"), context.Param("input"))
	if err != nil {
		msg := fmt.Sprintf("Error: %s", err.Error())
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, status)

}

func GetVolume(context echo.Context) error {
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "read-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

	status, err := di.GetVolume(context.Param("address"), context.Param("input"))
	if err != nil {
		msg := fmt.Sprintf("Error: %s", err.Error())
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, status)
}
