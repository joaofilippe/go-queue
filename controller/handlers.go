package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/joaofilippe/go-queue/models"
	redisclient "github.com/joaofilippe/go-queue/redis"
)


var (
	adminCode = 5654
	lastID    = 0
)

// EnterOnQueue insert the user on the queue and returns the render of the success screen
func EnterOnQueue(c echo.Context) error {
	ctx := context.Background()

	uF := new(models.UserForm)
	if err := c.Bind(uF); err != nil {
		return c.JSON(http.StatusBadRequest, err)

	}

	u := &models.User{
		Name:  uF.Name,
		Cpf:   uF.Cpf,
		Phone: uF.Phone,
	}

	if uF.Priority == "on" {
		u.Priority = true
	}

	if uF.Fever == "on" {
		u.Symptoms = append(u.Symptoms, "fever")
	}

	if uF.Nausea == "on" {
		u.Symptoms = append(u.Symptoms, "nausea")
	}

	if uF.Headache == "on" {
		u.Symptoms = append(u.Symptoms, "headache")
	}

	if uF.HeartDisease == "on" {
		u.Symptoms = append(u.Symptoms, "heart disease")
	}

	if uF.Air == "on" {
		u.Symptoms = append(u.Symptoms, "air")
	}

	if uF.Throat == "on" {
		u.Symptoms = append(u.Symptoms, "throat")
	}

	if uF.Hypertension == "on" {
		u.Symptoms = append(u.Symptoms, "hypertension")
	}

	if uF.LowPressure == "on" {
		u.Symptoms = append(u.Symptoms, "low pressure")
	}

	if uF.Vertigo == "on" {
		u.Symptoms = append(u.Symptoms, "vertigo")
	}

	if uF.Alergy == "on" {
		u.Symptoms = append(u.Symptoms, "alergy")
	}

	if uF.Diabetes == "on" {
		u.Symptoms = append(u.Symptoms, "diabetes")
	}

	if uF.Others != "" {
		u.Others = uF.Others
	}

	u.EnterOn = time.Now()

	if lastID == 1000 {
		lastID = 0
	}

	lastID++

	u.ID = lastID

	key := "standart"
	if u.Priority {
		key = "priority"
	}

	queue := redisclient.GetQueueFromRedis(ctx, key)
	queue.InsertNewUser(*u)
	redisclient.SendQueueToRedis(ctx, queue, key)

	queuePrior := redisclient.GetQueueFromRedis(ctx, "priority")
	queueStd := redisclient.GetQueueFromRedis(ctx, "standart")

	u.Place, _ = queuePrior.GetPlaceByID(u.ID)
	if u.Place == -1 {
		u.Place, _ = queueStd.GetPlaceByID(u.ID)
	}

	return RenderSuccessScreen(c, u)
}

// GetQueue gets the queue ordered by priority
func GetQueue(c echo.Context) error {
	ctx := context.Background()
	queueStd := redisclient.GetQueueFromRedis(ctx, "standart")
	queuePrior := redisclient.GetQueueFromRedis(ctx, "priority")

	queue := &models.UserQueue{}

	for _, u := range queuePrior.Queue {
		queue.Queue = append(queue.Queue, u)
	}

	for _, u := range queueStd.Queue {
		queue.Queue = append(queue.Queue, u)
	}

	message := struct {
		Data models.UserQueue `json:"data"`
	}{
		*queue,
	}

	return c.JSON(http.StatusFound, message)
}


// GetPlaceOnListByID is a function that returns the position of the user on the list
func GetPlaceOnListByID(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	queueStd := redisclient.GetQueueFromRedis(ctx, "standart")
	queuePrior := redisclient.GetQueueFromRedis(ctx, "priority")

	place, user := queuePrior.GetPlaceByID(id)
	if place == -1 {
		place, user = queueStd.GetPlaceByID(id)
		if place != -1 {
			place += queuePrior.Len()
		}
	}

	message := struct {
		Messsage string      `json:"message"`
		User     models.User `json:"user"`
		Place    int         `json:"place"`
	}{}

	if place == -1 {
		message.Messsage = "User not found"
		message.User = models.User{}
		message.Place = -1
		return c.JSON(http.StatusNotFound, message)
	}

	message.Messsage = "User found"
	message.User = user
	message.Place = place

	return c.JSON(http.StatusOK, message)
}



// CallNext calls the next user
func CallNext(c echo.Context) error {
	ctx := context.Background()
	code := c.Param("code")

	if code != "5654" {
		return c.JSON(http.StatusUnauthorized, struct {
			Message string `json:"message"`
		}{
			"Código inválido",
		})

	}

	key := "priority"
	queue := redisclient.GetQueueFromRedis(ctx, key)

	if queue.Len() == 0 {
		key = "standart"
		queue = redisclient.GetQueueFromRedis(ctx, key)
	}

	if queue.Len() == 0 {
		return c.JSON(http.StatusNoContent, struct {
			Message string `json:"message"`
		}{
			"empty queue",
		})

	}

	user := queue.Remove()

	redisclient.SendQueueToRedis(ctx, queue, key)

	return c.JSON(http.StatusOK, user)

}
