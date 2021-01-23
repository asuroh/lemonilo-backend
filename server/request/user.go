package request

// UserRequest ...
type UserRequest struct {
	Name     string `json:"name" `
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserUpdateRequest ...
type UserUpdateRequest struct {
	Name     string `json:"name" `
	UserName string `json:"user_name" validate:"required"`
}

// UserUploadImageRequest ...
type UserUploadImageRequest struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// UserLoginRequest ....
type UserLoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password" validate:"required"`
}
