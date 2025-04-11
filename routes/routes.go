package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	controller "github.com/wikasdude/blog-backend/controllers"
	"github.com/wikasdude/blog-backend/middleware"
)

func InitRoutes(userController *controller.UserController, postController *controller.PostController, categoryController *controller.CategoryController) {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userController.RegisterUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userController.LoginUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {

		// Extract ID from path
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 3 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		id := pathParts[2]
		fmt.Println("", r.Method)
		switch r.Method {
		case http.MethodGet:
			userController.GetUserByID(w, r, id)
		case http.MethodPut:
			userController.UpdateUser(w, r, id)
		case http.MethodDelete:
			userController.DeleteUser(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		//userController.GetUserByID(w, r, id)

	})
	http.HandleFunc("/api/blog-post", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postController.CreatePost(w, r)
			//middleware.AuthMiddleware(postController.CreatePost)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/api/blog-posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			postController.GetAllPosts(w, r)
			//postController.GetPaginatedPosts(w, r)
			return
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/blog-post/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("path is:", r.URL.Path)
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 4 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		fmt.Println("line no 62", len(pathParts))
		fmt.Println(pathParts[2])

		id, err := strconv.Atoi(pathParts[3])
		fmt.Println(id, " error is:", err)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Apply middleware only for PUT and DELETE methods
		if r.Method == http.MethodPatch || r.Method == http.MethodDelete {
			middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPatch:
					postController.UpdatePost(w, r)
				case http.MethodDelete:
					postController.DeletePost(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			})(w, r)
		} else {
			switch r.Method {
			case http.MethodGet:
				postController.GetPostByID(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			categoryController.CreateCategory(w, r)
		} else if r.Method == http.MethodGet {
			categoryController.GetAllCategories(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
