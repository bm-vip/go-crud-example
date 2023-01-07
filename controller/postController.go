package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go-crud-example/initializers"
	"go-crud-example/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func Signup(c *gin.Context) {
	//get the email/pass from request body
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// hash the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// create the user
	user := models.User{
		Email:    body.Email,
		Password: string(hashPassword),
	}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func Signin(c *gin.Context) {
	//get the email/pass from request body
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	//look up the user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	//compare pass of request body with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	//send the token
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

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
