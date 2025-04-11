// @title Blog API
// @version 1.0
// @description This is a simple blog backend using net/http in Go.
// @termsOfService http://swagger.io/terms/

// @contact.name Vikas
// @contact.email vikas82393@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes https
// @host blog-app-api-7t7q.onrender.com
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/wikasdude/blog-backend/config"
	controller "github.com/wikasdude/blog-backend/controllers"
	repository "github.com/wikasdude/blog-backend/repositories"
	service "github.com/wikasdude/blog-backend/services"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/wikasdude/blog-backend/docs"
	router "github.com/wikasdude/blog-backend/routes"
)

func init() {
	_ = godotenv.Load() // silently loads from .env if present
}

func main() {
	db, err := config.ConnectDB()
	fmt.Println(err)
	defer db.Close()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postController := controller.NewPostController(postService)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryController := controller.NewCategoryController(categoryService)

	router.InitRoutes(userController, postController, categoryController)
	log.Println("Server running on :8080")
	http.Handle("/swagger/", enableCORS(httpSwagger.WrapHandler))
	http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux))
}
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
