package repos

import "github.com/go-xorm/xorm"

type AuthRepo interface {
	CheckPassword(user, sign string) bool
}

type Auth struct {
	db *xorm.Engine
}

func NewAuth(db *xorm.Engine) *Auth {
	return &Auth{db: db}
}

func (a *Auth) CheckPassword(user, sign string) bool {
	return true
}
