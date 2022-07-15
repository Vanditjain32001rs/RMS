package models

type SubAdminModel struct {
	Name     string   `json:"name" db:"name"`
	Email    string   `json:"email" db:"email"`
	Username string   `json:"userName" db:"username"`
	Password string   `json:"password" db:"password"`
	Role     []string `json:"role" db:"role"`
}

//type FetchSubAdminModel struct {
//	ID       string   `json:"id"`
//	Name     string   `json:"name" db:"name"`
//	Email    string   `json:"email" db:"email"`
//	Username string   `json:"userName" db:"username"`
//	Role     []string `json:"role" db:"role"`
//}
