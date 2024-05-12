package controller

import (
	"context"
	"strconv"
	"text/template"

	"github.com/joaofilippe/go-queue/models"
	redisclient "github.com/joaofilippe/go-queue/redis"
	"github.com/labstack/echo/v4"
)

const tPath = "./screens/"

var tpl *template.Template

// LoadTemplates loads the templates
func LoadTemplates() {
	tpl = template.Must(template.ParseGlob(tPath + "*.html"))
}

// RenderHomeScreen render the home screen
func RenderHomeScreen(c echo.Context) error {
	return tpl.ExecuteTemplate(c.Response(), "home.html", nil)
}

// RenderSuccessScreen render the success screen
func RenderSuccessScreen(c echo.Context, u *models.User) error {
	return tpl.ExecuteTemplate(c.Response(), "success.html", u)
}

// RenderGetPlaceOnListByIDScreen is a function that returns the position of the user on the list
func RenderGetPlaceOnListByIDScreen(c echo.Context) error {
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

// RenderCallNextScreen render the call next screen
func RenderCallNextScreen(c echo.Context) error {
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
