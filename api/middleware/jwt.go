package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func JWTAuthentications(c *fiber.Ctx) error {
	fmt.Println("-- JWT auth --")

	tokenHeader := c.Get("X-API-Token")
	if tokenHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No authentication token provided",
		})
	}

	// Parse and validate the token
	err := parseToken(tokenHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid token: %v", err),
		})
	}

	return nil
}

func parseToken(Tokenstr string) error {

	token, err := jwt.Parse(Tokenstr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthoried")
		}

		secret := os.Getenv("JWT_SECRET")
		fmt.Println("NEVER EVER PRINT SECRET OR SHARE!!", secret)
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse token", err)
		return fmt.Errorf("unauthoried")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["nbf"])
	}

	return fmt.Errorf("unauthoried")
}
