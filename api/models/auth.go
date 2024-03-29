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
