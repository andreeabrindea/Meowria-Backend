package db

import (
	"time"
)

type Session struct {
	UserID       int
	SessionToken string
	Expiry       time.Time
}
