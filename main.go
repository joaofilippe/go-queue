package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	"github.com/joaofilippe/go-queue/models"
)

var queue = models.UserQueue{}
var adminCode = 456321

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	e := echo.New()
	e.GET("/", HelloWorld)
	e.GET("/queue", GetQueue)
	e.GET("/admin/:code", ValidateAdmin)
	e.POST("/enter", EnterOnQueue)

	e.Start(":8000")
}

func HelloWorld(c echo.Context) error {
	message := struct {
		Message string `json:"message"`
	}{
		"It's alive",
	}
	return c.JSON(http.StatusOK, message)
}

func EnterOnQueue(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	queue.InsertNewUser(*u)
	return c.JSON(http.StatusAccepted, u)
}

func GetQueue(c echo.Context) error {
	message := struct {
		Data models.UserQueue `json:"data"`
	}{
		queue,
	}

	return c.JSON(http.StatusFound, message)
}

func ValidateAdmin(c echo.Context) error {
	codeStr := c.Param("code")
	code, _ := strconv.Atoi(codeStr)
	if code == adminCode {
		return c.JSON(http.StatusOK,
			struct {
				Message string `json:"message"`
			}{
				"ok"})
	}

	return c.JSON(http.StatusUnauthorized, struct {
		Message string `json:"message"`
	}{
		Message: "not ok"})

}
