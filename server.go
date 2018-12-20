package main

import (
	"net/http"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/v2/auth"

	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/london-audio-microservice/handlers"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
)

func main() {
	port := ":8009"
	router := common.NewRouter()

	// Use the `secure` routing group to require authentication
	write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))
	router.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))
	router.GET("/status", echo.WrapHandler(http.HandlerFunc(health.Check)))

	read.GET("/raw", handlers.RawInfo)

	//functionality
	write.POST("/raw", handlers.Raw)
	write.GET("/:address/:input/volume/mute", handlers.Mute)
	write.GET("/:address/:input/volume/unmute", handlers.UnMute)
	write.GET("/:address/:input/volume/set/:level", handlers.SetVolume)

	//status
	read.GET("/:address/:input/volume/level", handlers.GetVolume)
	read.GET("/:address/:input/mute/status", handlers.GetMute)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
