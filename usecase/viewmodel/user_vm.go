package viewmodel

// UserVM ...
type UserVM struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
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
