package models

type SubAdminModel struct {
	Name     string   `json:"name" db:"name" validate:"required"`
	Email    string   `json:"email" db:"email" validate:"required & email"`
	Username string   `json:"userName" db:"username" validate:"required"`
	Password string   `json:"password" db:"password" validate:"required"`
	Role     []string `json:"role" db:"role" validate:"required"`
}

//type FetchSubAdminModel struct {
//	ID       string   `json:"id"`
//	Name     string   `json:"name" db:"name"`
//	Email    string   `json:"email" db:"email"`
//	Username string   `json:"userName" db:"username"`
//	Role     []string `json:"role" db:"role"`
//}
