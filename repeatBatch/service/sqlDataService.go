package service

import (
	"container/list"
	"log"
	"push_v2/module"

	"github.com/jmoiron/sqlx"
)

// SQLDataService ...
type SQLDataService struct {
	DbURL              string
	SqlLimitPushTarget int
	GoroutineCnt       int
	db                 *sqlx.DB
	rows               *sqlx.Rows
	err                error
}

// GetRepeatMessage ...
func (s *SQLDataService) GetRepeatMessage(listSlice []*list.List) {
	defer func() {
		s.rows.Close()
		module.DBClose(s.db)
	}()
	s.db = module.DBconn(s.DbURL)
	s.rows, s.err = s.db.Queryx(SelectPushBatchMsg, s.SqlLimitPushTarget)
	if s.err != nil {
		log.Printf("SelectPushBatchMsg : %s\n", s.err)
	}

	indexGoroutine := s.GoroutineCnt
	m := new(module.Message)
	for s.rows.Next() {
		err := s.rows.StructScan(&m)
		if err != nil {
			log.Printf("SelectPushBatchMsg StructScan : %s\n", err)
		}
		indexGoroutine--
		listSlice[indexGoroutine].PushBack(*m)
		if indexGoroutine == 0 {
			indexGoroutine = s.GoroutineCnt
		}
	}
}
