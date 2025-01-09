package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"hotelReservetion/types"
	"os"
)

func JWTAuthentications(c *fiber.Ctx) error {
	token := c.Cookies("access_token")

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No authentication token provided",
		})
	}

	// Parse and validate the token
	claims, err := ValidateToken(token, types.AccessToken)
	if err != nil {
		// If token is expired, check if we have a refresh token
		if err.Error() == "token expired" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired",
				"code":  "TOKEN_EXPIRED",
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid token: %v", err),
		})
	}

	if tokenType, ok := claims["type"].(string); !ok || types.TokenType(tokenType) != types.AccessToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token type",
		})
	}

	c.Locals("claims", claims)

	return c.Next()
}

func ValidateToken(tokenStr string, tokenType types.TokenType) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Choose secret based on token type
		var secret string
		if tokenType == types.AccessToken {
			secret = os.Getenv("JWT_SECRET")
		} else {
			secret = os.Getenv("JWT_REFRESH_SECRET")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// Verify token type
	if claimType, ok := claims["type"].(string); !ok || types.TokenType(claimType) != tokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}
