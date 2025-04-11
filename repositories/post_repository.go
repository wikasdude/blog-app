package repository

import (
	"database/sql"
	"fmt"

	"github.com/wikasdude/blog-backend/model"
)

type PostRepository struct {
	db *sql.DB
}

// Constructor
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

// Create a new post
func (r *PostRepository) CreatePost(post *model.Post) error {
	query := `INSERT INTO posts (user_id, title, description, category_id, body, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, post.UserID, post.Title, post.Description, post.CategoryID, post.Body).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

// Get a post by ID
func (r *PostRepository) GetPostByID(id int) (*model.Post, error) {
	query := `SELECT id, user_id, title, description, category_id, body, created_at, updated_at FROM posts WHERE id = $1`

	post := &model.Post{}
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Description,
		&post.CategoryID,
		&post.Body,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// Update a post
func (r *PostRepository) UpdatePost(post *model.Post) error {
	query := `UPDATE posts SET title = $1, description = $2, category_id = $3, body = $4, updated_at = NOW() WHERE id = $5`
	_, err := r.db.Exec(query, post.Title, post.Description, post.CategoryID, post.Body, post.ID)
	return err
}

// Delete a post
func (r *PostRepository) DeletePost(id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Get all posts (optional utility method)
func (r *PostRepository) GetAllPosts() ([]*model.Post, error) {
	query := `SELECT id, user_id, title, description, category_id, body, created_at, updated_at FROM posts ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Description,
			&post.CategoryID,
			&post.Body,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
func (r *PostRepository) GetPaginatedPosts(limit, offset int, sort string, order string, search string) ([]*model.Post, error) {
	query := fmt.Sprintf(`
	SELECT id, user_id, title, description, category_id, body, created_at, updated_at
	FROM posts
	WHERE title ILIKE '%%' || $3 || '%%'
	   OR description ILIKE '%%' || $3 || '%%'
	ORDER BY %s %s
	LIMIT $1 OFFSET $2
`, sort, order)
	rows, err := r.db.Query(query, limit, offset, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Description, &post.CategoryID, &post.Body, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
