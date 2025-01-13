package types

type PaginationOptions struct {
	Page     int64
	PageSize int64
	SortBy   string
	SortDesc bool
}

type Filter map[string]interface{}
