package models

import (
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Location struct {
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}

type Users struct {
	Name     string   `json:"name" db:"name"`
	Email    string   `json:"email" db:"email"`
	Username string   `json:"userName" db:"username"`
	Password string   `json:"password" db:"password"`
	Role     []string `json:"role" db:"role"`
	Location Location `json:"location"`
}

type UsersDetail struct {
	ID       string     `json:"id" db:"id"`
	Name     string     `json:"name" db:"name"`
	Email    string     `json:"email" db:"email"`
	Username string     `json:"userName" db:"username"`
	Role     []string   `json:"role" db:"role"`
	Location []Location `json:"location"`
}

type UpdateUsersModel struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Username string `json:"userName" db:"username"`
}

type UserFetchModel struct {
	ID         string   `json:"userID" db:"id"`
	Name       string   `json:"name" db:"name"`
	Email      string   `json:"email" db:"email"`
	Username   string   `json:"userName" db:"username"`
	Role       []string `json:"role" db:"role"`
	TotalCount int      `json:"-" db:"total_count"`
}

type UserFetch struct {
	TotalCount int              `json:"totalCount"`
	User       []UserFetchModel `json:"user"`
}

type UserFetchAdmin struct {
	TotalCount int           `json:"totalCount"`
	User       []UsersDetail `json:"user"`
}

type UserModel struct {
	Name     string   `json:"name" db:"name"`
	Email    string   `json:"email" db:"email"`
	Username string   `json:"userName" db:"username"`
	Password string   `json:"password" db:"password"`
	Location Location `json:"location"`
}
type UserLocation struct {
	Username string   `json:"userName"`
	UserLoc  Location `json:"userLocation"`
}

type UsersLocations struct {
	UserID    string  `db:"user_id"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}

type UserRoleID struct {
	UserID    string    `json:"userID" db:"id"`
	UserRole  string    `json:"userRole" db:"role"`
	CreatedBy uuid.UUID `json:"createdBy" db:"created_by"`
}

type RoleStruct struct {
	UserID   string `json:"userID" db:"user_id"`
	UserRole string `json:"userRole" db:"role"`
}
type ContextMap struct {
	UserID   string `json:"userID" db:"id"`
	Username string `json:"userName" db:"username"`
	UserRole string `json:"role" db:"role"`
}

type AddRoleModel struct {
	ID       string `json:"id" db:"user_id"`
	Username string `json:"userName" db:"username"`
	Role     string `json:"role" db:"role"`
}

type Pagination struct {
	TotalCount int              `json:"totalCount" db:"total_count"`
	User       []UserFetchModel `json:"usersDetail"`
}
