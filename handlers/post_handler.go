package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gopost-api/server"
	"github.com/gopost-api/services"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePostHandler(c *server.Context) {
	userID := c.GetUserID()
	if userID == 0 {
		RespondError(c.RWriter, NewAppError("Usuario no autenticado", http.StatusUnauthorized))
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.BindJSON(&req); err != nil {
		RespondError(c.RWriter, NewAppError("Datos inválidos", http.StatusBadRequest))
		return
	}

	post, err := h.postService.CreatePost(userID, req.Title, req.Content)
	if err != nil {
		RespondError(c.RWriter, NewAppError(err.Error(), http.StatusBadRequest))
		return
	}

	RespondJSON(c.RWriter, http.StatusCreated, map[string]interface{}{
		"message": "Post creado exitosamente",
		"post":    post,
	})
}

func (h *PostHandler) GetPostsHandler(c *server.Context) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		RespondError(c.RWriter, NewAppError(err.Error(), http.StatusInternalServerError))
		return
	}

	RespondJSON(c.RWriter, http.StatusOK, map[string]interface{}{
		"posts": posts,
	})
}

func (h *PostHandler) GetPostHandler(c *server.Context) {
	pathParts := strings.Split(c.Request.URL.Path, "/")
	if len(pathParts) < 3 {
		RespondError(c.RWriter, NewAppError("ID de post inválido", http.StatusBadRequest))
		return
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		RespondError(c.RWriter, NewAppError("ID de post inválido", http.StatusBadRequest))
		return
	}

	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		RespondError(c.RWriter, NewAppError("Post no encontrado", http.StatusNotFound))
		return
	}

	RespondJSON(c.RWriter, http.StatusOK, map[string]interface{}{
		"post": post,
	})
}

func (h *PostHandler) UpdatePostHandler(c *server.Context) {
	userID := c.GetUserID()
	if userID == 0 {
		RespondError(c.RWriter, NewAppError("Usuario no autenticado", http.StatusUnauthorized))
		return
	}

	pathParts := strings.Split(c.Request.URL.Path, "/")
	if len(pathParts) < 3 {
		RespondError(c.RWriter, NewAppError("ID de post inválido", http.StatusBadRequest))
		return
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		RespondError(c.RWriter, NewAppError("ID de post inválido", http.StatusBadRequest))
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.BindJSON(&req); err != nil {
		RespondError(c.RWriter, NewAppError("Datos inválidos", http.StatusBadRequest))
		return
	}

	post, err := h.postService.UpdatePost(uint(id), userID, req.Title, req.Content)
	if err != nil {
		RespondError(c.RWriter, NewAppError(err.Error(), http.StatusBadRequest))
		return
	}

	RespondJSON(c.RWriter, http.StatusOK, map[string]interface{}{
		"message": "Post actualizado exitosamente",
		"post":    post,
	})
}

func (h *PostHandler) DeletePostHandler(c *server.Context) {
	userID := c.GetUserID()
	if userID == 0 {
		RespondError(c.RWriter, NewAppError("Usuario no autenticado", http.StatusUnauthorized))
		return
	}

	pathParts := strings.Split(c.Request.URL.Path, "/")
	if len(pathParts) < 3 {
		RespondError(c.RWriter, NewAppError("ID de post inválido", http.StatusBadRequest))
		return
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		RespondError(c.RWriter, NewAppError("ID de post inválido", http.StatusBadRequest))
		return
	}

	if err := h.postService.DeletePost(uint(id), userID); err != nil {
		RespondError(c.RWriter, NewAppError(err.Error(), http.StatusBadRequest))
		return
	}

	RespondJSON(c.RWriter, http.StatusOK, map[string]interface{}{
		"message": "Post eliminado exitosamente",
	})
}

func (h *PostHandler) GetPostMeHandler(c *server.Context) {
	userID := c.GetUserID()
	if userID == 0 {
		RespondError(c.RWriter, NewAppError("Usuario no autenticado", http.StatusUnauthorized))
		return
	}

	posts, err := h.postService.GetPostsByUserID(userID)
	if err != nil {
		RespondError(c.RWriter, NewAppError(err.Error(), http.StatusInternalServerError))
		return
	}

	RespondJSON(c.RWriter, http.StatusOK, map[string]interface{}{
		"posts": posts,
	})
}
