package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservetion/types"
)

func GetUserFromClaims(claims jwt.MapClaims) (*types.User, error) {
	id, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'id' claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'email' claim")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'role' claim")
	}

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, fmt.Errorf("invalid 'id' firnat %v", err)
	}

	user := &types.User{
		ID:    objectID,
		Email: email,
		Role:  types.UserRole(role),
	}

	return user, nil
}
