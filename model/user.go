package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // 角色: admin, user, guest
	Email    string `json:"email"`
	Key      string `json:"key"`
	ExpireAt string `json:"expireAt"`
}
