package main

import (
	"context"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/joaofilippe/go-queue/controller"
	redisclient "github.com/joaofilippe/go-queue/redis"
)

func init() {
	controller.LoadTemplates()
}
func main() {
	ctx := context.Background()
	controller.LoadTemplates()
	redisclient.LoadRedisClient(ctx)

	e := echo.New()

	// Handlers
	e.POST("/enter", controller.EnterOnQueue)
	e.GET("/queue", controller.GetQueue)
	e.GET("/user/:id", controller.GetPlaceOnListByID)
	e.DELETE("/call/:code", controller.CallNext)
	
	// Renders
	e.GET("/", controller.RenderHomeScreen)
	e.GET("/place/:id", controller.RenderGetPlaceOnListByIDScreen)
	e.GET("/call", controller.RenderCallNextScreen)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	err := e.Start(":" +port)
	if err != nil {
		panic(err)
	}
}

