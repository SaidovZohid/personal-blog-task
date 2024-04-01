package models

type SignUpReq struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"email,required"`
	Code  string `json:"code" binding:"required"`
}

type LoginReq struct {
	Email      string `json:"email" binding:"email,required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me" default:"false"`
}

type UserData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
