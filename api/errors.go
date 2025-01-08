package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(*Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

type Error struct {
	Code int    `json:code`
	Msg  string `json:msg`
}

func (c *Error) Error() string {
	return c.Msg
}
func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func ErrInvalidID() *Error {
	return &Error{
		Code: http.StatusBadRequest,
		Msg:  "Invalid ID",
	}
}

func ErrBadRequest() *Error {
	return &Error{
		Code: http.StatusBadRequest,
		Msg:  "Invalid JSON request",
	}
}

func ErrNotFoundResources() *Error {
	return &Error{
		Code: http.StatusNotFound,
		Msg:  "Not found resources",
	}
}

func ErrUnauthorized() *Error {
	return &Error{
		Code: http.StatusUnauthorized,
		Msg:  "Unauthotized",
	}
}
