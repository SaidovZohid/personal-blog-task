package models

type Error struct {
	Code    int
	Error   string
	Message string
}

type ResponseSuccess struct {
	Code    int
	Message string
}

type UserLoginAndValidateReq struct {
	Info        UserInfo `json:"info"`
	AcceccToken string   `json:"accecc_token"`
}

type UserInfo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
