# Ejemplo de Uso del TimeoutMiddleware

## Uso Básico

Para usar el middleware de timeout en rutas específicas:

```go
// En cmd/api/main.go

import (
    "time"
    "github.com/gopost-api/middleware"
)

func main() {
    // ... código de inicialización ...

    // Aplicar timeout de 5 segundos a una ruta específica
    app.Get("/posts", middleware.TimeoutMiddleware(5*time.Second)(postHandler.GetPostsHandler))

    // O aplicar a una ruta protegida con múltiples middlewares
    app.Get("/posts/me",
        middleware.AuthMiddleware(
            middleware.TimeoutMiddleware(10*time.Second)(
                postHandler.GetPostMeHandler,
            ),
        ),
    )
}
```

## Ejemplo Completo

```go
package main

import (
    "log"
    "time"

    "github.com/gopost-api/config"
    "github.com/gopost-api/database"
    "github.com/gopost-api/handlers"
    "github.com/gopost-api/middleware"
    "github.com/gopost-api/repositories"
    "github.com/gopost-api/server"
    "github.com/gopost-api/services"
)

func main() {
    cfg := config.LoadConfig()

    if err := database.Connect(cfg.DatabaseURL); err != nil {
        log.Fatal("Error al conectar a la base de datos:", err)
    }
    defer database.Close()

    // Inicializar capas
    userRepo := repositories.NewUserRepository(database.DB)
    postRepo := repositories.NewPostRepository(database.DB)

    userService := services.NewUserService(userRepo)
    postService := services.NewPostService(postRepo)

    userHandler := handlers.NewUserHandler(userService)
    postHandler := handlers.NewPostHandler(postService)

    app := server.New()

    // Rutas sin timeout (usan el timeout por defecto del servidor: 15s)
    app.Get("/health", health)

    // Rutas de autenticación con timeout de 5 segundos
    app.Post("/auth/signup",
        middleware.TimeoutMiddleware(5*time.Second)(userHandler.SignUpHandler),
    )
    app.Post("/auth/login",
        middleware.TimeoutMiddleware(5*time.Second)(userHandler.LoginHandler),
    )

    // Rutas protegidas con timeout de 3 segundos
    app.Get("/auth/me",
        middleware.AuthMiddleware(
            middleware.TimeoutMiddleware(3*time.Second)(userHandler.MeHandler),
        ),
    )

    // Rutas de posts públicas con timeout de 10 segundos
    app.Get("/posts",
        middleware.TimeoutMiddleware(10*time.Second)(postHandler.GetPostsHandler),
    )
    app.Get("/posts/{id}",
        middleware.TimeoutMiddleware(5*time.Second)(postHandler.GetPostHandler),
    )

    // Rutas de posts protegidas con timeout de 7 segundos
    app.Post("/posts",
        middleware.AuthMiddleware(
            middleware.TimeoutMiddleware(7*time.Second)(postHandler.CreatePostHandler),
        ),
    )
    app.Put("/posts/{id}",
        middleware.AuthMiddleware(
            middleware.TimeoutMiddleware(7*time.Second)(postHandler.UpdatePostHandler),
        ),
    )
    app.Delete("/posts/{id}",
        middleware.AuthMiddleware(
            middleware.TimeoutMiddleware(5*time.Second)(postHandler.DeletePostHandler),
        ),
    )
    app.Get("/posts/me",
        middleware.AuthMiddleware(
            middleware.TimeoutMiddleware(10*time.Second)(postHandler.GetPostMeHandler),
        ),
    )

    if err := app.RunServer(cfg.Port); err != nil {
        log.Fatal("Error al iniciar el servidor:", err)
    }
}

func health(c *server.Context) {
    c.JSON(200, map[string]interface{}{
        "status":  "ok",
        "message": "El servicio está funcionando correctamente",
    })
}
```

## Timeouts Recomendados por Tipo de Operación

| Operación                  | Timeout Recomendado | Razón                             |
| -------------------------- | ------------------- | --------------------------------- |
| Health Check               | Sin timeout o 1s    | Operación muy rápida              |
| Login/SignUp               | 3-5s                | Operaciones de hash pueden tardar |
| Lectura simple (GET by ID) | 3-5s                | Query rápida con índice           |
| Listado (GET all)          | 10-15s              | Puede retornar muchos registros   |
| Creación (POST)            | 5-7s                | Insert + validaciones             |
| Actualización (PUT)        | 5-7s                | Update + validaciones             |
| Eliminación (DELETE)       | 3-5s                | Delete es rápido                  |

## Manejo de Errores de Timeout

Cuando se agota el timeout, el middleware responde automáticamente con:

```json
{
  "error": "La operación tardó demasiado tiempo"
}
```

Con status HTTP `408 Request Timeout`.

## Ventajas de Este Enfoque

1. **Protección contra queries lentas**: Evita que una query lenta bloquee recursos
2. **Mejor experiencia de usuario**: El cliente recibe una respuesta rápida en lugar de esperar indefinidamente
3. **Protección de recursos**: Libera conexiones de BD y memoria
4. **Configuración granular**: Cada ruta puede tener su propio timeout según sus necesidades
5. **Composable**: Se puede combinar con otros middlewares (auth, logging, etc.)

## Testing

Para probar el timeout, puedes crear un handler lento:

```go
func slowHandler(c *server.Context) {
    // Simular operación lenta
    time.Sleep(6 * time.Second)
    c.JSON(200, map[string]string{"message": "Esto no debería llegar"})
}

// En main.go
app.Get("/slow", middleware.TimeoutMiddleware(3*time.Second)(slowHandler))
// Al hacer GET /slow, debería retornar 408 después de 3 segundos
```
