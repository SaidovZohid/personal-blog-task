package models

type Error struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ResponseSuccess struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UserLoginAndValidateReq struct {
	Info        UserInfo `json:"info"`
	AcceccToken string   `json:"access_token"`
	RememberMe  bool     `json:"remember_me"`
}

type UserInfo struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
