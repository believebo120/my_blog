package main

import (
	"log"
	"my_blog/config"
	"my_blog/middleware"
	"my_blog/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// 初始化数据库连接
	config.InitDB()

	// 创建路由器
	router := mux.NewRouter()

	// 初始化路由
	routes.InitializeRoutes(router)

	// 应用CORS中间件
	routerWithCors := middleware.CorsMiddleware(router)

	// 启动服务器
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", routerWithCors); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
