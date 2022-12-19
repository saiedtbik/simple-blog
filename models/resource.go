package models

type Resource struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	CONTENT       string `json:"content" binding:"required"`
	USER_ID       uint   `json:"userId" binding:"required"`
	RESOURCE_TYPE uint   `json:"resourceType" binding:"required"`
	TITLE         string `json:"title"`
	REF_ID        uint   `json:"refId"`
}

type SearchResourceInput struct {
	ID           int    `json:"id"`
	InContent    string `json:"inContent"`
	UserId       uint   `json:"userId"`
	ResourceType uint   `json:"resourceType"`
	InTitle      string `json:"inTitle"`
	RefId        uint   `json:"refId"`
}
