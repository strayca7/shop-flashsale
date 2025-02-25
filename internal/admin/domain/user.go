package admindomain

import (
	"time"
)

type AdminUser struct {
	id string
	email string
	paswordHash string
	roles []Role
	createAt time.Time
}

func (u *AdminUser) ValidatePassword(password string) bool {
	return false
}
