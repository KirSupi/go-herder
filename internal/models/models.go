package models

type Session struct {
	ID        string
	UserAgent string
	IP        string
	CreatedAt int // unix timestamp
}
