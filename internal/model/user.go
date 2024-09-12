package model

import "time"

type At struct {
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type User struct {
	ID        string
	FirstName *string
	LastName  *string

	HashedRefreshToken *string

	At
}
