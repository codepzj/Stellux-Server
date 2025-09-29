package domain

import "time"

type File struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	FileName  string
	Url       string
	Dst       string
}
