package dto

// Standard Response

type SuccessResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
}

type ErrorResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
}

type PaginatedResponse struct {
    Success    bool        `json:"success"`
    Data       interface{} `json:"data"`
    Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
}