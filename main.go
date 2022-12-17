package main

import (
	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/handlers"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/Dhruv9449/mou/pkg/oauth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.Connect()
	models.InitializeModels()
	oauth.InitializeAuth()
	// driveutils.InitializeDrive()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Authorization",
		AllowOrigins:     "*, http://localhost:3000",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,DELETE,PATCH,PUT,OPTIONS",
	}))

	handlers.AuthRouter(app)
	handlers.UserRouter(app)
	handlers.BlogRouter(app)
	// handlers.DriveRouter(app)

	app.Listen(":8000")
}
