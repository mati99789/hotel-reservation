package api

import (
	"hotelReservetion/types"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

func AuthorizeRole(roles ...types.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get claims directly - they're set by JWT middleware
		claims := c.Locals("claims")
		if claims == nil {
			return &AuthorizationError{Message: "No authentication claims found"}
		}

		// Convert claims to the correct type
		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			return &AuthorizationError{Message: "Invalid token claims format"}
		}

		// Get role from claims
		roleInterface, exists := mapClaims["role"]
		if !exists {
			return &AuthorizationError{Message: "No role specified in token"}
		}

		userRole := types.UserRole(roleInterface.(string))

		// Check if user's role matches any of the allowed roles
		for _, allowedRole := range roles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		return &AuthorizationError{Message: "Insufficient privileges"}
	}
}
