package routes

import (
	"github.com/freeletics/lector/controllers"
	"github.com/gofiber/fiber/v2"
)

func ConfusableRoute(route fiber.Router) {
	route.Post("", controllers.CheckConfusable)
}
