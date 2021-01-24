package request

// UserRequest ...
type UserRequest struct {
	Name     string `json:"name" `
	Email    string `json:"email"`
	Address  string `json:"address"`
	UserName string `json:"username" validate:"required"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserUpdateRequest ...
type UserUpdateRequest struct {
	Name     string `json:"name" `
	UserName string `json:"username" validate:"required"`
}

// UserUpdatePasswordRequest ...
type UserUpdatePasswordRequest struct {
	Password string `json:"password" validate:"required"`
}

// UserUploadImageRequest ...
type UserUploadImageRequest struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// UserLoginRequest ....
type UserLoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password" validate:"required"`
}
