package request

type CreateCategoryRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Color     *string `json:"color"`
	SortOrder int     `json:"sortOrder"`
	IsDefault bool    `json:"isDefault"`
}

type UpdateCategoryRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Color     *string `json:"color"`
	SortOrder int     `json:"sortOrder"`
}