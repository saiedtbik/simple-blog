package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"withCasbin/component"
	"withCasbin/models"
)

func ReadResource(c *gin.Context) {
	// some stuff
	// blahblah...

	c.JSON(200, component.RestResponse{Code: 1, Message: "read resource successfully", Data: "resource"})
}

func WriteResource(c *gin.Context) {
	// some stuff
	// blahblah...

	c.JSON(200, component.RestResponse{Code: 1, Message: "write resource successfully", Data: "resource"})
}

func SearchResource(c *gin.Context) {
	var input models.SearchResourceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var res []models.Resource
	var query string

	query += "SELECT * FROM resources where 1=1"
	id := input.ID
	if  id>0 {
		query +=" and id = @ID"
		//args = append(args,  strconv.Itoa(id))
	/*	if err := models.DB.Where("id = ?",id).Find(&res).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
*/
	}

	userId := input.UserId
	if  userId>0 {
		query +=" and user_id = @UserId"
	//	args = append(args ,strconv.Itoa(int(userId)))

		/*if err := models.DB.Where("user_id = ?", userId).Find(&res).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}*/

	}

	refId := input.RefId
	if  refId>0 {
		query +=" and ref_id = @RefId"

	}

	resourceType := input.ResourceType
	if  resourceType>0 {
		query +=" and resource_type = @ResourceType"
	//	args = append(args , strconv.Itoa(int(resourceType)))

		/*if err := models.DB.Where("resource_type = ?", resourceType).Find(&res).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}*/

	}

	inTitle := input.InTitle
	if len(inTitle) > 0  {
		query +=" and title like @InTitle"
		//args = append(args , "'%"+inTitle+"%'")

		/*if err := models.DB.Where("title like ?", "%"+inTitle+"%").Find(&res).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}*/
	}

	inContent := input.InContent
	if len(inContent) > 0  {
		query +=" and content like @InContent"
		//args = append(args , "'%"+inContent+"%'")

		/*if err := models.DB.Where("content like ?","%"+inContent+"%").Find(&res).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}*/
	}
	if len(inContent) > 0 || len(inTitle) > 0 || resourceType>0 || userId>0 ||id>0 {
		if err :=models.DB.Raw(query,models.SearchResourceInput{ID:input.ID, ResourceType: input.ResourceType, UserId: input.UserId, RefId: input.RefId, InContent: "%"+input.InContent+"%", InTitle: "%"+input.InTitle+"%"}).Find(&res).Error; err != nil {
			//if err := models.DB.Where(query, args).Find(&res).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
	}else{
		if err :=models.DB.Find(&res).Error; err != nil {
			//if err := models.DB.Where(query, args).Find(&res).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func CreateResource(c *gin.Context){
	var input models.Resource
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create comment or post
	res := models.Resource{UserId: input.UserId, Content: input.Content, ResourceType: input.ResourceType, Title: input.Title}
	models.DB.Create(&res)

	c.JSON(http.StatusOK, gin.H{"data": res})
}