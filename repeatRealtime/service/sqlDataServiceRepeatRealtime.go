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

// GetRepeatRealtimeMessage ...
func (rs *SQLDataService) GetRepeatRealtimeMessage(listSlice []*list.List) {
	defer func() {
		rs.rows.Close()
		module.DBClose(rs.db)
		if r := recover(); r != nil {
			log.Printf("[Recover] GetRepeatRealtimeMessage() : %s\n", r)
		}
	}()
	rs.db = module.DBconn(rs.DbURL)
	rs.rows, rs.err = rs.db.Queryx(Select_RepeatRealtimeMsg, rs.SqlLimitPushTarget)
	if rs.err != nil {
		log.Printf("Select_RepeatRealtimeMsg : %s\n", rs.err)
	}

	indexGoroutine := rs.GoroutineCnt
	m := new(module.Message)
	for rs.rows.Next() {
		err := rs.rows.StructScan(&m)
		if err != nil {
			log.Printf("Select_RepeatRealtimeMsg StructScan : %s\n", err)
		}
		indexGoroutine--
		listSlice[indexGoroutine].PushBack(*m)
		if indexGoroutine == 0 {
			indexGoroutine = rs.GoroutineCnt
		}
	}
}
