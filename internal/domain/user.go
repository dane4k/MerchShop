package domain

type User struct {
	ID             int    `json:"id" db:"id"`
	Username       string `json:"name" db:"username"`
	PasswordHashed string `json:"password" db:"password"`
	Coins          int    `json:"coins" db:"balance"`
}
