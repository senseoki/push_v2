package service

import (
	"container/list"
	"log"
	"push_v2/module"
	"strings"

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

// GetRealtimeMessage ...
func (rs *SQLDataService) GetRealtimeMessage(listSlice []*list.List) {
	defer func() {
		rs.rows.Close()
		module.DBClose(rs.db)
		if r := recover(); r != nil {
			log.Printf("[Recover] GetRealtimeMessage() : %s\n", r)
		}
	}()
	rs.db = module.DBconn(rs.DbURL)
	rs.rows, rs.err = rs.db.Queryx(Select_PushTargetRealtime, rs.SqlLimitPushTarget)
	if rs.err != nil {
		log.Printf("Select_PushTargetRealtime : %s\n", rs.err)
	}

	indexGoroutine := rs.GoroutineCnt
	m := new(module.Message)
	for rs.rows.Next() {
		err := rs.rows.StructScan(&m)
		if err != nil {
			log.Printf("Select_PushTargetRealtime StructScan : %s\n", rs.err)
		}
		indexGoroutine--
		listSlice[indexGoroutine].PushBack(*m)
		if indexGoroutine == 0 {
			indexGoroutine = rs.GoroutineCnt
		}
	}
}

// InsertRealtimeStatus ...
func InsertRealtimeStatus(inVals []string, db *sqlx.DB) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recover] InsertRealtimeStatus() : %s\n", r)
		}
	}()
	_, err := db.Exec(Insert_PushTargetRealtimeStatus + strings.Join(inVals, ","))
	if err != nil {
		log.Printf("Insert_PushTargetRealtimeStatus StructScan : %s\n", err)
	}
}
