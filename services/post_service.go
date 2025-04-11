package service

import (
	"github.com/wikasdude/blog-backend/model"
	repository "github.com/wikasdude/blog-backend/repositories"
)

type PostService struct {
	postRepo *repository.PostRepository
}

// Constructor function to initialize PostService
func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

// Create a new post
func (s *PostService) CreatePost(post *model.Post) error {
	return s.postRepo.CreatePost(post)
}

// Get post by ID
func (s *PostService) GetPostByID(id int) (*model.Post, error) {
	return s.postRepo.GetPostByID(id)
}

// Update a post
func (s *PostService) UpdatePost(post *model.Post) error {
	return s.postRepo.UpdatePost(post)
}

// Delete a post
func (s *PostService) DeletePost(id int) error {
	return s.postRepo.DeletePost(id)
}

// List all posts (optional utility method)
func (s *PostService) GetAllPosts() ([]*model.Post, error) {
	return s.postRepo.GetAllPosts()
}
func (s *PostService) GetPaginatedPosts(page int, limit int, sort string, order string, search string) ([]*model.Post, error) {
	offset := (page - 1) * limit
	return s.postRepo.GetPaginatedPosts(limit, offset, sort, order, search)
}
