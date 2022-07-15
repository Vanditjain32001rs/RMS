package models

type Credentials struct {
	Username string `json:"userName" db:"username"`
	Password string `json:"password" db:"password"`
}
