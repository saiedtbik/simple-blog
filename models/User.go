package models

type User struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	UserName  string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Role struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	RoleName  string `json:"roleName" binding:"required"`
}

type CreateUserInput struct {
	UserName  string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateRoleInput struct {
	RoleName  string `json:"roleName" binding:"required"`
}

type UserRole struct{
	ID     uint   `json:"id" gorm:"primary_key"`
	User string   `json:"userName" binding:"required"`
	Role string   `json:"roleName" binding:"required"`

}

type RolePermissionResource struct{
	ID     uint   `json:"id" gorm:"primary_key"`
	Role string   `json:"roleName" binding:"required"`
	Permission string `json:"permission" binding:"required"`
	Resource string `json:"resource" binding:"required"`
}

