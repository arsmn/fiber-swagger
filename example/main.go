package main

import (
	fiberSwagger "github.com/arsmn/fiber-swagger"
	_ "github.com/arsmn/fiber-swagger/example/docs"
	"github.com/gofiber/fiber"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	app := fiber.New()

	app.Use(fiberSwagger.New())

	app.Listen(8080)
}
