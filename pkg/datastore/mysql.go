package datastore

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type MysqlConfig struct {
	Url         string
	MaxIdleConn int
	MaxOpenConn int
	LogMode     bool
}

func NewMysql(c MysqlConfig) (conn *xorm.Engine, err error) {
	log.Printf("%+v", c)
	conn, err = xorm.NewEngine("mysql", c.Url)
	if err != nil {
		return
	}

	if err = conn.DB().Ping(); err != nil {
		return
	}

	if c.MaxIdleConn == 0 {
		c.MaxIdleConn = 10
	}

	if c.MaxOpenConn == 0 {
		c.MaxOpenConn = 100
	}

	conn.DB().SetMaxIdleConns(c.MaxIdleConn)
	conn.DB().SetMaxOpenConns(c.MaxOpenConn)
	conn.ShowSQL(true)

	// FIXME temporary
	// err = conn.Sync(&models.CronJobs{})

	return
}
