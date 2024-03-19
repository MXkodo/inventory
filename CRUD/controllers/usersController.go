package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MXkodo/inventory/CRUD/initializers"
	"github.com/MXkodo/inventory/CRUD/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var body struct {
	FIO string
	UserName string
	Password string
}

func ExtractUserFIO(c *gin.Context) (string, error) {
	tokenString, err := c.Cookie("Authorization")
    if err != nil {
        return "", fmt.Errorf("authorization cookie is missing or invalid")
    }

	token, err := jwt.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("anvalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	userFIO, ok := claims["fio"].(string)
	if !ok {
		return "", fmt.Errorf("user FIO not found in token")
	}

	return userFIO, nil
}

func Signup(c *gin.Context) {

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	//Create User
	user := models.User{FIO: body.FIO, UserName: body.UserName, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	//Respond
	c.JSON(http.StatusOK, gin.H{})
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Login(c *gin.Context) {
	//Get the UserName and pass off req body
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//Look up requested user
	var user models.User
	initializers.DB.First(&user, "user_name = ?", body.UserName)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный логин или пароль",
		})
		return
	}
	//Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный логин или пароль",
		})
		return
	}
	//Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
    	"exp": time.Now().Add(time.Hour * 24).Unix(),
    	"fio": user.FIO,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to create token",
		})
		return
	}
	//send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*12, "", "", false, false)
	
	c.JSON(http.StatusOK, gin.H{
    "success": true,
    "fio": user.FIO,
})
}

func Logout(c *gin.Context){
	c.SetCookie("Authorization", "", -1, "","localhost", false, true)
	
		c.JSON(http.StatusOK, gin.H{
			"OK": "Cookie has been deleted",
		})
		
	
}
