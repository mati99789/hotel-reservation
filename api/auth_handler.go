package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"hotelReservetion/db"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var auth AuthParams
	if err := c.BodyParser(&auth); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), auth.Email)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credantials")
		}
		return err
	}

	fmt.Println("Authenticated", user)
	return nil
}
