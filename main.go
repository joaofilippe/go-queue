package main

import (
	"context"

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

	e.Start(":8000")
}

// criar um id ao invés de usar o CPF na params
// criar a página para o usuário ver sua posição na fila
// inserir na página de chamar uma lista para o funcionário ver os próximos 5 da fila
// inserir mais informações para a triagem
// o que se encaixa no prioritário (pcd, gestante, recém-nascido e idoso)

// whish: possibilitar que o funcionário altere o grau de prioridade do usuário
