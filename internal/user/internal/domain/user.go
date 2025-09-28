package domain

import (
	"time"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  string
	Nickname  string
	RoleId    int
	Avatar    string
	Email     string
}

type Page struct {
	PageNo   int64
	PageSize int64
}
