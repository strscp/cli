package models

// PaginatedResponse represents a Laravel paginated API response.
type PaginatedResponse[T any] struct {
	Data  []T             `json:"data"`
	Links PaginationLinks `json:"links"`
	Meta  PaginationMeta  `json:"meta"`
}

// PaginationLinks contains pagination navigation URLs.
type PaginationLinks struct {
	First string  `json:"first"`
	Last  string  `json:"last"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}

// PaginationMeta contains pagination metadata.
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	From        int `json:"from"`
	To          int `json:"to"`
}
