package models

type Credentials struct {
	Username string `json:"userName" db:"username" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
}
