package main

import (
	"context"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/joaofilippe/go-queue/handlers"
	redisclient "github.com/joaofilippe/go-queue/redis"
)

func init() {
	handlers.LoadTemplates()
}
func main() {
	ctx := context.Background()
	handlers.LoadTemplates()
	redisclient.LoadRedisClient(ctx)

	e := echo.New()
	e.GET("/", handlers.HomeScreen)
	e.POST("/enter", handlers.EnterOnQueue)
	e.GET("/queue", handlers.GetQueue)
	e.GET("/user/:id", handlers.GetPlaceOnListByID)
	e.GET("/place/:id", handlers.GetPlaceOnListByIDScreen)
	e.GET("/call", handlers.CallNextScreen)
	e.GET("/see", handlers.SeePatients)
	e.DELETE("/call/:code", handlers.CallNext)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	err := e.Start(":" +port)
	if err != nil {
		panic(err)
	}
}

