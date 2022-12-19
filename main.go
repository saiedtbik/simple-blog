package main

import (
	"fmt"
	filedApter "github.com/casbin/casbin/v2/persist/file-adapter"
	//	sqldApter "github.com/casbin/casbin/v2/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"withCasbin/models"

	//"withCasbin/component"
	"withCasbin/handler"
	"withCasbin/middleware"
)

var (
	router *gin.Engine
)

func init() {
	models.ConnectDatabase()

	// Initialize  casbin adapter
	adapter := filedApter.NewAdapter("config/basic_policy.csv")

	// Initialize gin router
	router = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig)) // CORS configuraion

	router.POST("/user/login", handler.Login)
	router.POST("/user", handler.CreateUser)
	router.POST("/role", handler.CreateRole)
	router.POST("/user/assign-role", handler.CreateUserRole)
	router.POST("/rule", handler.CreateRule)

	resource := router.Group("/api")
	resource.Use(middleware.Authenticate())
	{
		resource.GET("/post", middleware.Authorize("post", "read", adapter), handler.SearchResource)
		resource.POST("/comment", middleware.Authorize("comment", "write", adapter), handler.CreateResource)
		resource.POST("/post", middleware.Authorize("post", "write", adapter), handler.CreateResource)
		resource.GET("/comment", middleware.Authorize("comment", "read", adapter), handler.SearchResource)
	}
}

func main() {
/*	defer func() {
		err := component.DB.Close()
		if err != nil {
			log.Println(fmt.Errorf("failed to close DB connection: %w", err))
		}
	}()
*/
	err := router.Run(":8081")
	if err != nil {
		log.Fatalln(fmt.Errorf("faild to start Gin application: %w", err))
	}
}
