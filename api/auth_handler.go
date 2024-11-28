package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservetion/db"
	"hotelReservetion/types"
	"os"
	"time"
)

const (
	INVALIDCREDENTIALS = "Invalid Credentials"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(genericResp{
		Type: "error",
		Msg:  INVALIDCREDENTIALS,
	})
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
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, auth.Password) {
		return invalidCredentials(c)
	}

	token := createTokenFromUser(user)
	return c.JSON(&AuthResponse{
		User:  user,
		Token: token,
	})
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	validTill := now.Add(time.Hour * 4)
	claims := jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
		"exp":   validTill.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Error signing token")
	}

	return tokenString

}
