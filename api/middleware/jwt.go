package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func JWTAuthentications(c *fiber.Ctx) error {
	tokenHeader := c.Get("X-API-Token")
	if tokenHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No authentication token provided",
		})
	}

	// Parse and validate the token
	claims, err := validateToken(tokenHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid token: %v", err),
		})
	}

	fmt.Println("Claims: ", claims)

	return c.Next()
}

func validateToken(Tokenstr string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(Tokenstr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthoried")
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse token", err)
		return nil, fmt.Errorf("unauthoried")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}
	}

	return claims, nil
}
