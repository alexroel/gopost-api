package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gopost-api/handlers"
	"github.com/gopost-api/server"
)

// TimeoutMiddleware agrega un timeout específico a las peticiones
// Si la petición tarda más del timeout especificado, se cancela automáticamente
func TimeoutMiddleware(timeout time.Duration) func(server.HandleFunc) server.HandleFunc {
	return func(next server.HandleFunc) server.HandleFunc {
		return func(c *server.Context) {
			// Crear un context con timeout
			ctx, cancel := context.WithTimeout(c.Ctx, timeout)
			defer cancel()

			// Actualizar el context en el Context personalizado
			c.Ctx = ctx

			// Canal para saber cuándo terminó el handler
			done := make(chan bool, 1)

			// Ejecutar el handler en una goroutine
			go func() {
				next(c)
				done <- true
			}()

			// Esperar a que termine o se agote el timeout
			select {
			case <-done:
				// El handler terminó correctamente
				return
			case <-ctx.Done():
				// Se agotó el timeout
				handlers.RespondError(c.RWriter, handlers.NewAppError(
					"La operación tardó demasiado tiempo",
					http.StatusRequestTimeout,
				))
				return
			}
		}
	}
}
