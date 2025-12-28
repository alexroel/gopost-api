package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gopost-api/config"
	"github.com/gopost-api/handlers"
	"github.com/gopost-api/server"
)

func AuthMiddleware(next server.HandleFunc) server.HandleFunc {
	return func(c *server.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			handlers.RespondError(c.RWriter, handlers.NewAppError("Token no proporcionado", http.StatusUnauthorized))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handlers.RespondError(c.RWriter, handlers.NewAppError("Formato de token inválido", http.StatusUnauthorized))
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			handlers.RespondError(c.RWriter, handlers.NewAppError("Token inválido o expirado", http.StatusUnauthorized))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handlers.RespondError(c.RWriter, handlers.NewAppError("Claims inválidos", http.StatusUnauthorized))
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			handlers.RespondError(c.RWriter, handlers.NewAppError("User ID no encontrado en el token", http.StatusUnauthorized))
			return
		}

		c.SetUserID(uint(userID))

		next(c)
	}
}
