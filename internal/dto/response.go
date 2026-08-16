package dto

type Response[T any] struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}
