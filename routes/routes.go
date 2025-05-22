package routes

import (
	"my_blog/controllers"
	"my_blog/middleware"
	"my_blog/utils"

	"github.com/gorilla/mux"

	"github.com/gin-gonic/gin"
)

// InitializeRoutes 初始化路由
func InitializeRoutes(router *mux.Router) {
	// 创建控制器实例
	articleController := controllers.NewArticleController()
	userController := controllers.NewUserController()
	commentController := controllers.NewCommentController()
	categoryController := controllers.NewCategoryController()
	r := gin.Default()
	// 映射 URL 路径 `/uploads/` 到本地目录 `./uploads`（与你的目录结构一致）
	r.Static("/uploads", "./uploads")
	// 公共API，无需认证
	router.HandleFunc("/register", userController.Register).Methods("POST")
	router.HandleFunc("/login", userController.Login).Methods("POST")
	router.HandleFunc("/articles", articleController.GetArticles).Methods("GET")
	router.HandleFunc("/articles/{id}", articleController.GetArticle).Methods("GET")
	router.HandleFunc("/categories", categoryController.GetCategories).Methods("GET")
	router.HandleFunc("/categories/{id}", categoryController.GetCategory).Methods("GET")
	router.HandleFunc("/categories/{id}/articles", articleController.GetArticlesByCategory).Methods("GET")

	// 需要认证的API
	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	// 用户相关API
	authRouter.HandleFunc("/users/me", userController.GetCurrentUser).Methods("GET")
	authRouter.HandleFunc("/users/{id}/background-image", userController.UpdateUserBackgroundImage).Methods("POST")
	authRouter.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT")

	// 文章相关API
	authRouter.HandleFunc("/articles", articleController.CreateArticle).Methods("POST")
	authRouter.HandleFunc("/articles/{id}", articleController.UpdateArticle).Methods("PUT")
	authRouter.HandleFunc("/articles/{id}", articleController.DeleteArticle).Methods("DELETE")

	// 评论相关API
	authRouter.HandleFunc("/articles/{id}/comments", commentController.GetCommentsByArticle).Methods("GET")
	authRouter.HandleFunc("/articles/{id}/comments", commentController.CreateComment).Methods("POST")
	authRouter.HandleFunc("/comments/{id}", commentController.UpdateComment).Methods("PUT")
	authRouter.HandleFunc("/comments/{id}", commentController.DeleteComment).Methods("DELETE")

	// 需要管理员权限的API
	adminRouter := authRouter.PathPrefix("").Subrouter()
	adminRouter.Use(middleware.RoleMiddleware(utils.RoleAdmin))

	// 用户管理API
	adminRouter.HandleFunc("/users", userController.GetAllUsers).Methods("GET")
	adminRouter.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")
	adminRouter.HandleFunc("/users/{id}/role", userController.UpdateUserRole).Methods("PUT")

	// 分类管理API
	adminRouter.HandleFunc("/categories", categoryController.CreateCategory).Methods("POST")
	adminRouter.HandleFunc("/categories/{id}", categoryController.UpdateCategory).Methods("PUT")
	adminRouter.HandleFunc("/categories/{id}", categoryController.DeleteCategory).Methods("DELETE")
}
