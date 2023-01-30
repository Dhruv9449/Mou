package main

import (
	"log"
	"os"

	"github.com/Dhruv9449/mou/management"
	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/handlers"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/Dhruv9449/mou/pkg/oauth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/urfave/cli"
)

func main() {
	database.Connect()
	models.InitializeModels()
	oauth.InitializeAuth()
	// driveutils.InitializeDrive()

	// CLI
	var cli_app = cli.NewApp()

	cli_app.Name = "Mou"
	cli_app.Usage = "A simple blogging platform"
	cli_app.Version = "0.0.1"
	cli_app.Author = "Dhruv Shah"

	cli_app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Start the server",
			Action: func(c *cli.Context) error {
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
				return nil
			},
		},
		{
			Name:   "createsuperuser",
			Usage:  "Create a superuser",
			Action: management.CreateSuperuser,
		},
		{
			Name:   "viewallusers",
			Usage:  "View all users",
			Action: management.ViewAllUsers,
		},
	}

	err := cli_app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
