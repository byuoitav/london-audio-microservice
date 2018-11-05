package main

import (
	"net/http"

	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/london-audio-microservice/handlers"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := ":8009"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	//secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))
	router.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))
	router.GET("/status", echo.WrapHandler(http.HandlerFunc(health.Check)))

	router.GET("/raw", handlers.RawInfo)

	//functionality
	router.POST("/raw", handlers.Raw)
	router.GET("/:address/:input/volume/mute", handlers.Mute)
	router.GET("/:address/:input/volume/unmute", handlers.UnMute)
	router.GET("/:address/:input/volume/set/:level", handlers.SetVolume)

	//status
	router.GET("/:address/:input/volume/level", handlers.GetVolume)
	router.GET("/:address/:input/mute/status", handlers.GetMute)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
