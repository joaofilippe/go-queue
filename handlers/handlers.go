package handlers

import (
	"context"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/joaofilippe/go-queue/models"
	redisclient "github.com/joaofilippe/go-queue/redis"
)

const tPath = "./templates/"

var (
	tpl       *template.Template
	adminCode = 5654
	lastID    = 0
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

	return tpl.ExecuteTemplate(c.Response(), "success.html", u)
}

// GetQueue gets the queue
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

// GetPlaceOnListByCPF is a function that returns the position of the user on the list
func GetPlaceOnListByCPF(c echo.Context) error {
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

// GetPlaceOnListByIDScreen is a function that returns the position of the user on the list
func GetPlaceOnListByIDScreen(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	queueStd := redisclient.GetQueueFromRedis(ctx, "standart")
	queuePrior := redisclient.GetQueueFromRedis(ctx, "priority")
	queue := new(models.UserQueue)

	for _, u := range queuePrior.Queue {
		queue.Queue = append(queue.Queue, u)
	}

	for _, u := range queueStd.Queue {
		queue.Queue = append(queue.Queue, u)
	}

	place, user := queue.GetPlaceByID(id)
	user.Place = place

	return tpl.ExecuteTemplate(c.Response(), "place.html", &user)
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

func CallNextScreen(c echo.Context) error {
	ctx := context.Background()

	queuePrior := redisclient.GetQueueFromRedis(ctx, "priority")
	queueStd := redisclient.GetQueueFromRedis(ctx, "standart")

	if queuePrior.Len() == 0 && queueStd.Len() == 0 {
		return tpl.ExecuteTemplate(c.Response(), "empty.html", nil)

	}

	queue := []models.User{}
	for _, u := range queuePrior.Queue {
		queue = append(queue, u)
	}

	for _, u := range queueStd.Queue {
		queue = append(queue, u)
	}

	user := queue[0]

	user.Next = queue[1:]
	return tpl.ExecuteTemplate(c.Response(), "call.html", &user)

}

func SeePatients(c echo.Context) error {
	ctx := context.Background()
	queue := redisclient.GetQueueFromRedis(ctx, "standart")

	symptoms := []string{}

	for _, u := range queue.Queue {
		for _, s := range u.Symptoms {
			if s == "fever" {
				symptoms = append(symptoms, "febre")
			}

			if s == "nausea" {
				symptoms = append(symptoms, "náusea")
			}

			if s == "headache" {
				symptoms = append(symptoms, "dor de cabeça")
			}

			if s == "heart disease" {
				symptoms = append(symptoms, "doença cardíaca")
			}

			if s == "air" {
				symptoms = append(symptoms, "falta de ar")
			}

			if s == "throat" {
				symptoms = append(symptoms, "dor de garganta")
			}

			if s == "hypertension" {
				symptoms = append(symptoms, "hipertensão")
			}

			if s == "low pressure" {
				symptoms = append(symptoms, "pressão baixa")
			}

			if s == "vertigo" {
				symptoms = append(symptoms, "vertigem")
			}

			if s == "alergy" {
				symptoms = append(symptoms, "alergia")
			}

			if s == "diabetes" {
				symptoms = append(symptoms, "diabetes")
			}
		}
	}

	return tpl.ExecuteTemplate(c.Response(), "request.html", &queue)
}
