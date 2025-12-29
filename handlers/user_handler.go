package handlers

import (
	"net/http"

	"github.com/gopost-api/server"
	"github.com/gopost-api/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) SignUpHandler(c *server.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		RespondError(c.RWriter, NewAppError("Datos inválidos", http.StatusBadRequest))
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		RespondError(c.RWriter, NewAppError("Todos los campos son requeridos", http.StatusBadRequest))
		return
	}

	user, err := h.userService.SignUp(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		RespondError(c.RWriter, NewAppError(err.Error(), http.StatusBadRequest))
		return
	}

	RespondJSON(c.RWriter, http.StatusCreated, map[string]interface{}{
		"message": "Usuario registrado exitosamente",
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func (h *UserHandler) LoginHandler(c *server.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		RespondError(c.RWriter, NewAppError("Datos inválidos", http.StatusBadRequest))
		return
	}

	if req.Email == "" || req.Password == "" {
		RespondError(c.RWriter, NewAppError("Email y contraseña son requeridos", http.StatusBadRequest))
		return
	}

	token, err := h.userService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		RespondError(c.RWriter, NewAppError(err.Error(), http.StatusUnauthorized))
		return
	}

	RespondJSON(c.RWriter, http.StatusOK, map[string]interface{}{
		"message": "Inicio de sesión exitoso",
		"token":   token,
	})
}

func (h *UserHandler) MeHandler(c *server.Context) {
	userID := c.GetUserID()
	if userID == 0 {
		RespondError(c.RWriter, NewAppError("Usuario no autenticado", http.StatusUnauthorized))
		return
	}

	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		RespondError(c.RWriter, NewAppError("Usuario no encontrado", http.StatusNotFound))
		return
	}

	RespondJSON(c.RWriter, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
