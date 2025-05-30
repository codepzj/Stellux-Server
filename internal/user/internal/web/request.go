package web

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	RoleId   *int   `json:"role_id" binding:"required"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

type UpdateUserRequest struct {
	ID       string `json:"id" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UpdatePasswordRequest struct {
	ID       string `json:"id" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
