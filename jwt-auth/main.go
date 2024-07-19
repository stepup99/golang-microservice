package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("your_secret_key")

func main() {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		// Example user data, replace with actual user validation logic
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Validate user credentials (replace with actual validation logic)
		if username != "exampleUser" || password != "examplePass" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Example user ID and role, replace with actual user data
		userID := "exampleUserID"
		userRole := "admin" // Example role

		// Create JWT Token with user ID and role
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"role":    userRole,
			"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expiration (72 hours)
		})

		// Sign JWT Token
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			fmt.Println("Error signing token:", err) // Debugging line
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Return Token in Response
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})

	authorized := r.Group("/")
	authorized.Use(JWTAuthMiddleWare())

	authorized.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{"user_id": userID, "role": role})
	})

	adminGrp := authorized.Group("/admin", roleAuthMiddleWare("admin"))
	adminGrp.GET("/dashboard", func(c *gin.Context) {
		fmt.Println("welcome to the dashboard")
		c.JSON(http.StatusOK, gin.H{"success": "welcome to the dashboard "})
	})
	r.Run(":9091")
}

func roleAuthMiddleWare(roleType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exist := c.Get("role")
		if !exist || role != roleType {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient privilages"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Parse and Validate JWT Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println(tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Extract Claims and Set Context
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}

}
