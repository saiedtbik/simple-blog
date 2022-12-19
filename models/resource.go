package models

type Resource struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Content  string `json:"content" binding:"required"`
	UserId uint `json:"userId" binding:"required"`
	ResourceType uint `json:"resourceType" binding:"required"`
	Title string `json:"title"`
	RefId uint  `json:"refId"`
}

type SearchResourceInput struct {
	ID     int   `json:"id"`
	InContent  string `json:"inContent"`
	UserId uint `json:"userId"`
	ResourceType uint `json:"resourceType"`
	InTitle string `json:"inTitle"`
	RefId uint  `json:"refId"`
}
