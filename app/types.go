package app

import (
	"time"
)

type User struct {
	Id            uint64
	Username      string
	Email         string
	EncryptedPass string
	LastAccess    time.Time
	Verified      bool
}

type Budget struct {
	Id     uint64
	UserId uint64
	Income uint64
	Rent   uint64
	Wealth uint64
}
