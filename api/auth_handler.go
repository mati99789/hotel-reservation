package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservetion/db"
	"hotelReservetion/types"
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
	User *types.User `json:"user"`
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

	tokens, err := types.CreateTokenPair(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	accessCookie := fiber.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		MaxAge:   4 * 60 * 60,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	}

	c.Cookie(&accessCookie)

	refreshCookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/refresh",
		MaxAge:   7 * 24 * 60 * 60,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	}

	c.Cookie(&refreshCookie)

	return c.JSON(&AuthResponse{
		User: user,
	})
}

func (h *AuthHandler) HandleRefresh(c *fiber.Ctx) error {
	refreshToken := c.Params("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No refresh token provided",
		})
	}

	claims, err := ValidateToken(refreshToken, types.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	tokens, err := types.CreateTokenPair(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate tokens",
		})
	}

	// Set new access token cookie
	accessCookie := fiber.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		MaxAge:   4 * 60 * 60,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	}
	c.Cookie(&accessCookie)

	// Set new refresh token cookie
	refreshCookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/refresh",
		MaxAge:   7 * 24 * 60 * 60,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	}
	c.Cookie(&refreshCookie)

	return c.JSON(&AuthResponse{
		User: user,
	})
}
