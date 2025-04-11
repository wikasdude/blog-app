package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wikasdude/blog-backend/config"
	"github.com/wikasdude/blog-backend/model"
	service "github.com/wikasdude/blog-backend/services"
	"github.com/wikasdude/blog-backend/utils"
)

type PostController struct {
	postService *service.PostService
}

// Constructor
func NewPostController(postService *service.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// Create Post
// CreatePost godoc
// @Summary Create a blog post
// @Description Create a new post by authenticated user
// @Tags posts
// @Accept json
// @Produce json
// @Param post body model.CreatePostRequest true "Post Data"
// @Success 201 {object} model.Post
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/blog-post [post]
func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	var err error

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	tokenString := r.Header.Get("Authorization")
	claims, err := utils.ValidateJWT(strings.TrimPrefix(tokenString, "Bearer "))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	post.UserID = claims.UserID
	if post.Title == "" || post.Description == "" || post.UserID == 0 || post.CategoryID == 0 {
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(map[string]interface{}{
		// 	"status":  false,
		// 	"message": "All fields (title, description, user_id, category_id) are required",
		// })
		utils.SendError(w, http.StatusBadRequest, "All fields (title, description, user_id, category_id) are required", err)
		return
	}

	err = c.postService.CreatePost(&post)
	if err != nil {
		fmt.Println("line no 35:", err)
		//http.Error(w, "Failed to create post", http.StatusInternalServerError)
		utils.SendError(w, http.StatusInternalServerError, "Failed to create post", err)
		return
	}

	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(post)
	utils.SendSuccess(w, http.StatusOK, "Posts created successfully", post)

}

// func (c *PostController) GetAllPosts(w http.ResponseWriter, r *http.Request) {
// 	posts, err := c.postService.GetAllPosts()
// 	if err != nil {
// 		//http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
// 		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch posts", err)
// 		return
// 	}

// 	// w.WriteHeader(http.StatusOK)
// 	// json.NewEncoder(w).Encode(posts)
// 	utils.SendSuccess(w, http.StatusOK, "Posts fetched successfully", posts)
// }

// GetPostByID godoc
// @Summary Get a blog post by ID
// @Description Fetch a single post by its ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} model.Post
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Router /api/blog-post/{id} [get]
func (c *PostController) GetPostByID(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	fmt.Println("line no 62", len(pathParts))
	fmt.Println(pathParts[2])

	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Post not found", err)

	}

	post, err := c.postService.GetPostByID(id)
	if err != nil {
		// w.WriteHeader(http.StatusNotFound)
		// json.NewEncoder(w).Encode(map[string]interface{}{
		// 	"status":  false,
		// 	"message": fmt.Sprintf("Post not found with id: %d", id),
		// })
		utils.SendError(w, http.StatusNotFound, fmt.Sprintf("Post not found with id: %d", id), err)

		return
		// http.Error(w, "Post not found", http.StatusNotFound)
		// return
	}
	utils.SendSuccess(w, http.StatusOK, "Post fetched successfully", post)

}

// Update Post
// UpdatePost godoc
// @Summary Update an existing blog post
// @Description Allows the owner or an admin to update a blog post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body model.UpdatePostRequest true "Updated Post JSON"
// @Success 200 {object} model.Post
// @Failure 400 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/blog-post/{id} [patch]
func (c *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var err error

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.SendError(w, http.StatusBadRequest, "Invalid URL", err)
		return
	}

	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	// userID := utils.GetUserIDFromContext(r.Context())
	// role := utils.GetUserRoleFromContext(r.Context())
	tokenString := r.Header.Get("Authorization")
	claims, err := utils.ValidateJWT(strings.TrimPrefix(tokenString, "Bearer "))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID
	role := claims.Role

	existingPost, err := c.postService.GetPostByID(id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Post not found", err)
		return
	}

	if existingPost.UserID != userID && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "You are not allowed to update this post", nil)
		return
	}

	var updatedPost model.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	updatedPost.ID = id
	updatedPost.UserID = existingPost.UserID // keep original owner

	err = c.postService.UpdatePost(&updatedPost)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to update post", err)
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Post updated successfully", updatedPost)
}

// Delete Post
// Decode updated data
// DeletePost godoc
// @Summary Delete a blog post
// @Description Allows the owner or an admin to delete a blog post
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.APIResponse
// @Failure 403 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security BearerAuth
// @Router /api/blog-post/{id} [delete]
func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		utils.SendError(w, http.StatusBadRequest, "Invalid URL", "expected format: /api/blog-post/{id}")
		return
	}

	id, err := strconv.Atoi(pathParts[3])
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	tokenString := r.Header.Get("Authorization")
	claims, err := utils.ValidateJWT(strings.TrimPrefix(tokenString, "Bearer "))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := claims.UserID
	role := claims.Role

	post, err := c.postService.GetPostByID(id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Post not found", err.Error())
		return
	}

	fmt.Println(post.UserID, "  line no 201", userID)

	if post.UserID != userID && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "You are not allowed to delete this post", nil)
		return
	}

	err = c.postService.DeletePost(id)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to delete post", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, "Post deleted successfully", nil)
}

// GetAllPosts godoc
// @Summary Get all blog posts
// @Description Fetch all posts from the database
// @Tags posts
// @Produce json
// @Success 200 {array} model.Post
// @Failure 500 {object} utils.APIResponse
// @Router /api/blog-posts [get]
func (pc *PostController) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	if pageStr != "" && limitStr != "" {
		// Paginated
		page, err1 := strconv.Atoi(pageStr)
		limit, err2 := strconv.Atoi(limitStr)

		if err1 != nil || err2 != nil || page < 1 || limit < 1 {
			utils.SendError(w, http.StatusBadRequest, "Invalid pagination params", nil)
			return
		}
		sort := r.URL.Query().Get("sort")
		order := r.URL.Query().Get("order")
		if sort == "" {
			sort = config.DefaultSort
		}
		if order == "" {
			order = config.DefaultOrder
		}
		search := r.URL.Query().Get("search")
		//validSortFields := map[string]bool{"title": true, "created_at": true}
		if !config.ValidSortFields[sort] {
			sort = config.DefaultSort
		}
		if order != "asc" && order != "desc" {
			order = config.DefaultOrder
		}

		posts, err := pc.postService.GetPaginatedPosts(page, limit, sort, order, search)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch paginated posts", err)
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Paginated posts fetched", posts)
		return
	}

	// No pagination
	posts, err := pc.postService.GetAllPosts()
	if err != nil {
		fmt.Println("line no 317", err)
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch posts", err)
		return
	}

	utils.SendSuccess(w, http.StatusOK, "All posts fetched", posts)
}
