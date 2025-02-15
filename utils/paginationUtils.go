package utils

import (
	"hotelReservetion/types"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ExtractPaginationFromRequest(c *fiber.Ctx) *types.PaginationOptions {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Query("pageSize", "10"), 10, 64)
	sortBy := c.Query("sortBy", "")
	sortDesc := c.Query("sortDesc", "false") == "true"

	return &types.PaginationOptions{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}
}

type ResourceResponse[T any] struct {
	Result int `json:"results"`
	Data   []T `json:"data"`
	Page   int `jsons:"page"`
}

func NewResourceResponse[T any](result, page int, data []T) *ResourceResponse[T] {
	return &ResourceResponse[T]{
		Result: result,
		Data:   data,
		Page:   page,
	}
}
