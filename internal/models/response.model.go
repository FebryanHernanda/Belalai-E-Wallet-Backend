package models

type Response struct {
	IsSuccess bool   `json:"is_success"  example:"true"`
	Code      int    `json:"code,omitempty"  example:"200"`
	Page      int    `json:"page,omitempty"  example:"1"`
	Msg       string `json:"message,omitempty"  example:"Example message success..."`
}

type ErrorResponse struct {
	Response
	Err string `json:"error" example:"Error message..."`
}

type ResponseData struct {
	Response
	Data interface{} `json:"data,omitempty"`
}

// ===================== { Example error response for swagger, not use } =====================

type BadRequestResponse struct {
	IsSuccess bool   `json:"is_success" example:"false"`
	Code      int    `json:"code,omitempty" example:"400"`
	Err       string `json:"error" example:"Example bad request error..."`
}

type UnauthorizedResponse struct {
	IsSuccess bool   `json:"is_success" example:"false"`
	Code      int    `json:"code,omitempty" example:"401"`
	Err       string `json:"error" example:"Unauthorized : please login again..."`
}

type NotFoundResponse struct {
	IsSuccess bool   `json:"is_success" example:"false"`
	Code      int    `json:"code,omitempty" example:"404"`
	Err       string `json:"error" example:"Not Found"`
}

type InternalErrorResponse struct {
	IsSuccess bool   `json:"is_success" example:"false"`
	Code      int    `json:"code,omitempty" example:"500"`
	Err       string `json:"error" example:"Example Internal server error..."`
}
