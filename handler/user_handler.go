package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"withCasbin/component"
	"withCasbin/models"
)

func Login(c *gin.Context) {

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username, password := input.UserName,input.Password
	// If user has logged in, force him to log out firstly
	for iter := component.GlobalCache.Iterator(); iter.SetNext(); {
		info, err := iter.Value()
		if err != nil {
			continue
		}
		if string(info.Value()) == username {
			component.GlobalCache.Delete(info.Key())
			log.Printf("forced %s to log out\n", username)
			break
		}
	}
	//var users []models.User
	var user models.User
	// Struct
	models.DB.Where(&models.User{UserName: username, Password: password}).First(&user)
	if user.ID>0 {
		log.Println(fmt.Sprintf("%s has logged in.", username))
	}else {
		c.JSON(200, component.RestResponse{Message: "no such account"})
		return
	}


	// Generate random session id
	u, err := uuid.NewRandom()
	if err != nil {
		log.Println(fmt.Errorf("failed to generate UUID: %w", err))
	}
	sessionId := fmt.Sprintf("%s-%s", u.String(), username)
	// Store current subject in cache
	err = component.GlobalCache.Set(sessionId, []byte(username))
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to store current subject in cache: %w", err))
		return
	}
	// Send session id back to client in cookie
	c.SetCookie("current_subject", sessionId, 30*60, "/api", "", false, true)
	c.JSON(200, component.RestResponse{Code: 1, Message: username + " logged in successfully"})
}

func FindUserRole(username string) (roleName string){  // Get model if exist
	var userRole models.UserRole
	models.DB.Where(&models.UserRole{User: username}).First(&userRole)

	/*
		if err := models.DB.Where("user = ?", username).First(&user).Error; err != nil {
			log.Println(fmt.Sprintf("%s not found", username))
		}*/
	return userRole.Role
}

func CreateUser(c *gin.Context){
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user := models.User{UserName: input.UserName, Password: input.Password}
	models.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func CreateRole(c *gin.Context){
	var input models.CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create role
	role := models.Role{RoleName: input.RoleName}
	models.DB.Create(&role)

	c.JSON(http.StatusOK, gin.H{"data": role})
}

func CreateUserRole(c *gin.Context){
	var input models.UserRole
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create roleUser
	userRole := models.UserRole{Role: input.Role, User: input.User}
	models.DB.Create(&userRole)

	c.JSON(http.StatusOK, gin.H{"data": userRole})
}

func CreateRule(c *gin.Context){
	var input models.RolePermissionResource
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create rule
	rule := models.RolePermissionResource{Role: input.Role, Permission: input.Permission, Resource: input.Resource}
	models.DB.Create(&rule)

	c.JSON(http.StatusOK, gin.H{"data": rule})
}
