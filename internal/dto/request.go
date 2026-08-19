package dto

type RegisterRequest struct {
	Username string `json:"username" form:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

type CategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type RecordRequest struct {
	CategoryID  uint    `json:"category_id" validate:"required"`
	RecordType  string  `json:"record_type" validate:"required,oneof=income expense"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description" validate:"max=500"`
	RecordDate  string  `json:"record_date"`
}

type ExportReportRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type SummaryRequest struct {
	Month string `json:"month" validate:"required"` // MM-YYYY
}
