package service

import (
	"github.com/wikasdude/blog-backend/model"
	repository "github.com/wikasdude/blog-backend/repositories"
)

type CategoryService struct {
	Repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{Repo: repo}
}

func (s *CategoryService) CreateCategory(category *model.Category) error {
	return s.Repo.CreateCategory(category)
}

func (s *CategoryService) GetAllCategories() ([]*model.Category, error) {
	return s.Repo.GetAllCategories()
}
