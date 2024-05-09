package handlers

import (
	"context"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/joaofilippe/go-queue/models"
	"github.com/joaofilippe/go-queue/redis"
)

const tPath = "./templates/"

var (
	tpl       *template.Template
	adminCode = 5654
)

func LoadTemplates() {
	tpl = template.Must(template.ParseGlob(tPath + "*.html"))
}

// Screens
// HomeScreen is the home screen
func HomeScreen(c echo.Context) error {
	return tpl.ExecuteTemplate(c.Response(), "home.html", nil)
}

// EnterOnQueue enter on queue
func EnterOnQueue(c echo.Context) error {
	ctx := context.Background()

	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, err)

	}

	// u.EnterOn = time.Now()
	// u.Name = c.FormValue("name")
	// u.Cpf = c.FormValue("cpf")
	// u.Phone = c.FormValue("phone")
	// if c.FormValue("priority") == "on" {
	// 	u.Priority = true
	// }

	queue := new(models.UserQueue)
	if u.Priority == "on" {
		queue = redisclient.GetQueueFromRedis(ctx, "priority")
	} else {
		queue = redisclient.GetQueueFromRedis(ctx, "standart")
	}

	queue.InsertNewUser(*u)
	redisclient.SendQueueToRedis(ctx, queue)

	return c.JSON(http.StatusAccepted, u)
}

// GetQueue gets the queue
func GetQueue(c echo.Context) error {
	ctx := context.Background()
	queue := redisclient.GetQueueFromRedis(ctx, "standart")
	message := struct {
		Data models.UserQueue `json:"data"`
	}{
		*queue,
	}

	return c.JSON(http.StatusFound, message)
}

// GetPlaceOnList is a function that returns the position of the user on the list
func GetPlaceOnList(c echo.Context) error {
	ctx := context.Background()
	cpf := c.Param("cpf")
	queue := redisclient.GetQueueFromRedis(ctx, "standart")
	place := queue.GetPlaceByCPF(cpf)

	message := struct {
		Messsage string `json:"message"`
		Place    int    `json:"place"`
	}{}

	if place == -1 {
		message.Messsage = "User not found"
		message.Place = -1
		return c.JSON(http.StatusNotFound, message)
	}

	message.Messsage = "User found"
	message.Place = place

	return c.JSON(http.StatusOK, message)
}

// ValidateAdmin validate an admin
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

// CallNext calls the next user
func CallNext(c echo.Context) error {
	exitTime := time.Now()
	ctx := context.Background()
	queue := redisclient.GetQueueFromRedis(ctx, "standart")
	user := queue.Remove()

	timeOnQueue := exitTime.Sub(user.EnterOn)
	println("Time on queue: ", timeOnQueue)

	redisclient.SendQueueToRedis(ctx, queue)

	return c.JSON(http.StatusOK, user)

}
