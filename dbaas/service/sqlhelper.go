package service

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"strings"
)

type SqlHelper interface {
	And(sql string, args ...interface{})
	Or(sql string, args ...interface{})
	Raw(sql string, args ...interface{})
	Session(engine *xorm.Engine) *xorm.Session
}

type sqlHelper struct {
	sql  string
	args []interface{}
}

func (s *sqlHelper) setArgs(args []interface{}) {
	if len(args) != 0 {
		s.args = append(s.args, args...)
	}

}

func (s *sqlHelper) Raw(sql string, args ...interface{}) {
	s.sql = fmt.Sprintf("%v %v", s.sql, sql)
	s.setArgs(args)
}

func (s *sqlHelper) Session(engine *xorm.Engine) *xorm.Session {
	return engine.SQL(s.sql, s.args...)
}

func (s *sqlHelper) Or(sql string, args ...interface{}) {
	s.sql = fmt.Sprintf("%v or %v", s.sql, sql)
	s.setArgs(args)
}

func (s *sqlHelper) And(sql string, args ...interface{}) {
	s.sql = fmt.Sprintf("%v and %v", s.sql, sql)
	s.setArgs(args)
}

var _ SqlHelper = (*sqlHelper)(nil)

func NewSql(sql string, args ...interface{}) SqlHelper {
	if strings.LastIndex(sql, "where") >= 0 && strings.LastIndex(sql, "WHERE") >= 0 {
		sql = fmt.Sprintf("%v where 1=1", sql)
	}
	return &sqlHelper{sql: sql, args: args}
}
