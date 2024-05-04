package main

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	"github.com/joaofilippe/go-queue/models"
)

var (
	adminCode = 456321
	ctx       = context.Background()
	rCli      = &redis.Client{}
	tpl       *template.Template
)

const tPath = "./templates/"

func init() {
	tpl = template.Must(template.ParseGlob(tPath + "*.html"))
}
func main() {
	rCli = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	qStd := models.UserQueue{}
	qPri := models.UserQueue{}
	jsonQStd, err := json.Marshal(qStd)
	if err != nil {
		println(err)
	}
	jsonQPri, err := json.Marshal(qPri)
	if err != nil {
		println(err)
	}

	rCli.Set(ctx, "standart", string(jsonQStd), 0)
	rCli.Set(ctx, "priority", string(jsonQPri), 0)

	e := echo.New()
	e.GET("/", HomeScreen)
	e.GET("/queue", GetQueue)
	e.GET("/admin/:code", ValidateAdmin)
	e.POST("/enter", EnterOnQueue)
	e.DELETE("/call", CallNext)

	e.Start(":8000")
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

	u.EnterOn = time.Now()
	u.Name = c.FormValue("name")
	u.Cpf = c.FormValue("cpf")
	u.Phone = c.FormValue("phone")
	if c.FormValue("priority") == "on" {
		u.Priority = true
	}

	queue := new(models.UserQueue)
	if u.Priority {
		queue = GetQueueFromRedis(ctx, "priority")
	} else {
		queue = GetQueueFromRedis(ctx, "standart")
	}

	queue.InsertNewUser(*u)
	SendQueueToRedis(ctx, queue)

	return c.JSON(http.StatusAccepted, u)
}

// GetQueue gets the queue
func GetQueue(c echo.Context) error {
	ctx := context.Background()
	queue := GetQueueFromRedis(ctx, "standart")
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
	queue := GetQueueFromRedis(ctx, "standart")
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
	queue := GetQueueFromRedis(ctx, "standart")
	user := queue.Remove()

	timeOnQueue := exitTime.Sub(user.EnterOn)
	println("Time on queue: ", timeOnQueue)
	
	SendQueueToRedis(ctx, queue)

	return c.JSON(http.StatusOK, user)

}

// Redis funcs

// GetQueueFromRedis gets the queue from redis-cli
func GetQueueFromRedis(ctx context.Context, key string) *models.UserQueue {
	result := rCli.Get(ctx, key)
	queueStr, err := result.Result()
	if err != nil {
		println(err)
	}

	queue := new(models.UserQueue)

	json.Unmarshal([]byte(queueStr), queue)

	return queue
}

// SendQueueToRedis sends the new queue to Redis
func SendQueueToRedis(ctx context.Context, queue *models.UserQueue) {
	res, _ := json.Marshal(queue)
	rCli.Set(ctx, "standart", string(res), 0)
}
