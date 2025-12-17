package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/vn-fin/oms/docs"
	"github.com/vn-fin/oms/internal/config"
	"github.com/vn-fin/oms/internal/db"
	_ "github.com/vn-fin/oms/internal/models"
	"github.com/vn-fin/oms/internal/remote"
	v1Routes "github.com/vn-fin/oms/internal/v1/routes"
	v2Routes "github.com/vn-fin/oms/internal/v2/routes"

	"github.com/rs/zerolog/log"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	var err error

	// Init config
	if err = config.InitConfig(); err != nil {
		panic(err)
	}

	// Init Postgresql
	if err = db.InitPostgres(); err != nil {
		panic(err)
	}
	if err = db.InitPostgresUserDB(); err != nil {
		panic(err)
	}
	defer db.ClosePostgres()

	// Init Auth gRPC Client
	if err = remote.InitAuthGrpcClient(); err != nil {
		panic(err)
	}
	defer remote.CloseAuthGrpcClient()

	// Initialize the Fiber app
	app := fiber.New(fiber.Config{
		// If behind a proxy, use the X-Forwarded-For header
		// to get the client's real IP address.
		ProxyHeader: fiber.HeaderXForwardedFor,
		// Enable IP validation
		EnableIPValidation: true,
	})
	// Enable CORS for all origins
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Add swagger docs
	app.Get(fmt.Sprintf("/%s/swagger_docs/*", config.ServiceName), swagger.HandlerDefault)

	// Register routes for version 1
	v1Routes.SetupRoutes(app)
	v2Routes.SetupRoutes(app)

	// Start server
	err = app.Listen(fmt.Sprintf("%s:%d", "0.0.0.0", 3000))
	if err != nil {
		log.Error().Msgf("Error starting server: %v", err)
	} else {
		// Start the server
		log.Info().Msgf("Server running at http://%s:%d\n", "0.0.0.0", 3000)
	}
}
