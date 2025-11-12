package auth

type User struct {
	ID string `json:"id" binding:"required"`
	Token Tokens `json:"token" binding:"required"`
}
