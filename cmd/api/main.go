package main

import (
	"log"

	"github.com/gopost-api/config"
	"github.com/gopost-api/database"
	"github.com/gopost-api/handlers"
	"github.com/gopost-api/middleware"
	"github.com/gopost-api/repositories"
	"github.com/gopost-api/server"
	"github.com/gopost-api/services"
)

func health(c *server.Context) {
	c.JSON(200, map[string]interface{}{
		"status":  "ok",
		"message": "El servicio está funcionando correctamente",
	})
}

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Conectar a la base de datos
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	defer database.Close()

	// Inicializar repositorios
	userRepo := repositories.NewUserRepository(database.DB)
	postRepo := repositories.NewPostRepository(database.DB)

	// Inicializar servicios
	userService := services.NewUserService(userRepo)
	postService := services.NewPostService(postRepo)

	// Inicializar handlers
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)

	// Crear aplicación
	app := server.New()

	// Ruta de bienvenida
	app.Get("/health", health)

	// Rutas públicas - Usuarios
	app.Post("/auth/signup", userHandler.SignUpHandler)
	app.Post("/auth/login", userHandler.LoginHandler)

	// Rutas protegidas - Usuarios
	app.Get("/auth/me", middleware.AuthMiddleware(userHandler.MeHandler))

	// Rutas públicas - Posts
	app.Get("/posts", postHandler.GetPostsHandler)
	app.Get("/posts/{id}", postHandler.GetPostHandler)

	// Rutas protegidas - Posts
	app.Post("/posts", middleware.AuthMiddleware(postHandler.CreatePostHandler))
	app.Put("/posts/{id}", middleware.AuthMiddleware(postHandler.UpdatePostHandler))
	app.Delete("/posts/{id}", middleware.AuthMiddleware(postHandler.DeletePostHandler))
	app.Get("/posts/me", middleware.AuthMiddleware(postHandler.GetPostMeHandler))

	// Iniciar servidor
	if err := app.RunServer(); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
