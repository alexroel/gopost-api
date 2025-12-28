package services

import (
	"fmt"

	"github.com/gopost-api/models"
	"github.com/gopost-api/repositories"
)

type PostService struct {
	repo *repositories.PostRepository
}

func NewPostService(repo *repositories.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(userID uint, title, content string) (*models.Post, error) {
	if title == "" {
		return nil, fmt.Errorf("el título es requerido")
	}
	if content == "" {
		return nil, fmt.Errorf("el contenido es requerido")
	}

	post := &models.Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	if err := s.repo.Create(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	return s.repo.FindAll()
}

func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	return s.repo.FindByID(id)
}

func (s *PostService) GetPostsByUserID(userID uint) ([]models.Post, error) {
	return s.repo.FindByUserID(userID)
}

func (s *PostService) UpdatePost(postID, userID uint, title, content string) (*models.Post, error) {
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return nil, err
	}

	if post.UserID != userID {
		return nil, fmt.Errorf("no tienes permiso para actualizar este post")
	}

	if title == "" {
		return nil, fmt.Errorf("el título es requerido")
	}
	if content == "" {
		return nil, fmt.Errorf("el contenido es requerido")
	}

	post.Title = title
	post.Content = content

	if err := s.repo.Update(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(postID, userID uint) error {
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return fmt.Errorf("no tienes permiso para eliminar este post")
	}

	return s.repo.Delete(postID)
}
