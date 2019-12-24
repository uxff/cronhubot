package checker

import (
	"github.com/go-xorm/xorm"
)

type postgres struct {
	url string
}

func NewPostgres(url string) *postgres {
	return &postgres{url}
}

func (p *postgres) IsAlive() bool {
	conn, err := xorm.NewEngine("postgres", p.url)
	if err != nil {
		return false
	}
	defer conn.Close()

	if err = conn.DB().Ping(); err != nil {
		return false
	}

	return true
}
