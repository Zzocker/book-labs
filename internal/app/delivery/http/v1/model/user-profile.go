package model

type CreateProfileRequest struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Name     string `form:"name"`
	Password string `form:"password"`
}
