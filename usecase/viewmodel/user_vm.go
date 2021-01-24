package viewmodel

// UserVM ...
type UserVM struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	ImagePath string `json:"image_path"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

// UserUploadImageVM ...
type UserUploadImageVM struct {
	ID        string `json:"id"`
	Path      string `json:"path"`
	CreatedAt string `json:"created_at"`
}
