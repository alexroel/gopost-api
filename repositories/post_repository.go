package repositories

import (
	"database/sql"
	"fmt"

	"github.com/gopost-api/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
	query := "INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, post.UserID, post.Title, post.Content)
	if err != nil {
		return fmt.Errorf("error al crear post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error al obtener ID: %w", err)
	}

	post.ID = uint(id)
	return nil
}

func (r *PostRepository) FindAll() ([]models.Post, error) {
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) FindByID(id uint) (*models.Post, error) {
	post := &models.Post{}
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE id = ?"
	
	err := r.db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post no encontrado")
		}
		return nil, fmt.Errorf("error al buscar post: %w", err)
	}

	return post, nil
}

func (r *PostRepository) FindByUserID(userID uint) ([]models.Post, error) {
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE user_id = ? ORDER BY created_at DESC"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener posts del usuario: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error al escanear post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	query := "UPDATE posts SET title = ?, content = ? WHERE id = ?"
	result, err := r.db.Exec(query, post.Title, post.Content, post.ID)
	if err != nil {
		return fmt.Errorf("error al actualizar post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post no encontrado")
	}

	return nil
}

func (r *PostRepository) Delete(id uint) error {
	query := "DELETE FROM posts WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar eliminación: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post no encontrado")
	}

	return nil
}
