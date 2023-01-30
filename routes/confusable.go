package routes

import (
	"github.com/hingstarne/lector/controllers"
	"github.com/gofiber/fiber/v2"
)

func ConfusableRoute(route fiber.Router) {
	route.Post("", controllers.CheckConfusable)
}
