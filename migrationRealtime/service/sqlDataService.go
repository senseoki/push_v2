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
	DbURL string
	db    *sqlx.DB
	rows  *sqlx.Rows
	tx    *sqlx.Tx
	err   error
}

// GetMigrationRealtimeMessage ...
func (sds *SQLDataService) GetMigrationRealtimeMessage() *list.List {
	defer func() {
		sds.rows.Close()
		module.DBClose(sds.db)
		if r := recover(); r != nil {
			log.Printf("[Recover] GetMigrationRealtimeMessage() : %s\n", r)
		}
	}()

	sds.db = module.DBconn(sds.DbURL)
	sds.rows, sds.err = sds.db.Queryx(SelectPushTargetRealtimeStatus)
	if sds.err != nil {
		log.Printf("Select_PushTargetRealtimeStatus err : %s\n", sds.err)
	}

	m := new(module.Message)
	li := list.New()
	for sds.rows.Next() {
		err := sds.rows.StructScan(&m)
		if err != nil {
			log.Printf("Select_PushTargetRealtimeStatus StructScan err : %s\n", sds.err)
		}
		li.PushBack(*m)
	}
	return li
}

// ExecMigrationRealtime ...
func (sds *SQLDataService) ExecMigrationRealtime(li *list.List) {
	defer func() {
		sds.rows.Close()
		sds.tx.Commit()
		module.DBClose(sds.db)
		if r := recover(); r != nil {
			sds.tx.Rollback()
			log.Printf("[Recover] ExecMigrationRealtime() : %s\n", r)
		}
	}()

	sds.db = module.DBconn(sds.DbURL)

	values := make([]string, 0, 0)
	for e := li.Front(); e != nil; e = e.Next() {
		values = append(values, e.Value.(module.Message).PushTargetSeq)
	}

	sds.db.MustExec(InsertPushTargetRealtimeStatusLog + "(" + strings.Join(values, ",") + ")")
	sds.db.MustExec(InsertPushTargetRealtimeLog + "(" + strings.Join(values, ",") + ")")

	// InnoDB 는 트랜잭션 처리한다.
	sds.tx = sds.db.MustBegin()
	sds.tx.MustExec(DeletePushTargetRealtimeStatus + "(" + strings.Join(values, ",") + ")")
	sds.tx.MustExec(DeletePushTargetRealtime + "(" + strings.Join(values, ",") + ")")
}
