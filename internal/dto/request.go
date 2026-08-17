package dto

type RegisterRequest struct {
	Username string `json:"username" form:"username" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type RecordRequest struct {
	CategoryID  uint    `json:"category_id"`
	RecordType  string  `json:"record_type" validate:"required"`
	Amount      float64 `json:"amount" validate:"required"`
	Description string  `json:"description"`
	RecordDate  string  `json:"record_date"`
}
