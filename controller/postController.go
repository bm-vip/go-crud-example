package controller

import (
	"github.com/gin-gonic/gin"
	"go-crud-example/initializers"
	"go-crud-example/models"
)

func Create(c *gin.Context) {
	var body struct {
		Body  string
		Title string
	}
	c.Bind(&body)

	post := models.Post{
		Title: body.Title,
		Body:  body.Body,
	}
	result := initializers.DB.Create(&post)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"post": post,
	})
}
func GetAll(c *gin.Context) {
	var posts []models.Post
	initializers.DB.Find(&posts)
	c.JSON(200, gin.H{
		"posts": posts,
	})
}
func FindById(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	initializers.DB.Find(&post, id)
	c.JSON(200, gin.H{
		"post": post,
	})
}

func Update(c *gin.Context) {
	var body struct {
		Body  string
		Title string
	}
	c.Bind(&body)

	id := c.Param("id")
	var post models.Post
	initializers.DB.Find(&post, id)

	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})

	c.JSON(200, gin.H{
		"post": post,
	})
}

func DeleteById(c *gin.Context) {
	id := c.Param("id")

	initializers.DB.Delete(id)
	c.Status(200)
}
