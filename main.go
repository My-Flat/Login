package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"my-flat-login/internal/config"
	"my-flat-login/internal/controller"
	"my-flat-login/internal/middleware"
	"my-flat-login/internal/repository"
	"my-flat-login/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configurations
	dbCfg, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Failed to load database configuration: %v", err)
	}

	db, err := config.ConnectDB(dbCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get the underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	defer sqlDB.Close() // Close the database connection when the application exits

	// 3. Initialize Firebase
	ctx := context.Background()
	firebaseAuth, err := config.InitializeFirebase(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// 4. Set up repositories
	userRepository := repository.NewUserRepository(db)

	// 5. Set up services
	authService := service.NewAuthService(firebaseAuth, userRepository)

	// 6. Set up controllers
	authController := controller.NewAuthController(authService)

	// 7. Set up the Gin router
	router := gin.Default()

	// 8. Define API routes
	api := router.Group("/api")
	{
		api.POST("/login", authController.Login)

		// Protected routes
		protected := api.Group("/protected")
		{
			// Use the JWT middleware to protect this route
			jwtSecret := []byte("your-jwt-secret") // Replace with your actual JWT secret
			protected.Use(middleware.JWTAuthMiddleware(jwtSecret))

			protected.GET("/test", func(ctx *gin.Context) {
				// Access the user from the context (set by the middleware)
				user, _ := ctx.Get("user")
				ctx.JSON(http.StatusOK, gin.H{"message": "Protected route accessed", "user": user})
			})
		}
	}

	// 9. Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
