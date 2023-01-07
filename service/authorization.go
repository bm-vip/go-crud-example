package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go-crud-example/initializers"
	"go-crud-example/models"
	"net/http"
	"os"
	"strings"
	"time"
)

func ValidateToken(c *gin.Context) {
	fmt.Println("start validating token")
	//Get the token cookie of request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token is required",
			})
		}
		tokenString = strings.Replace(header, "bearer", "", 1)
		tokenString = strings.Replace(header, "Bearer", "", 1)
	}
	//Validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token is expired",
			})
		}
		//find the user with token sub
		var user = models.User{}
		initializers.DB.Find(&user, claims["sub"])
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token Unauthorized",
			})
		}
		//attach user to req
		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
