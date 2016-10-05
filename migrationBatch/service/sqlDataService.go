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

// GetMigrationBatchTarget ...
func (sds *SQLDataService) GetMigrationBatchTarget() *list.List {
	defer func() {
		sds.rows.Close()
		module.DBClose(sds.db)
		if r := recover(); r != nil {
			log.Printf("[Recover] GetMigrationBatchMessage() : %s\n", r)
		}
	}()

	sds.db = module.DBconn(sds.DbURL)
	sds.rows, sds.err = sds.db.Queryx(SelectPushTargetStatus)
	if sds.err != nil {
		log.Printf("SelectPushTargetStatus err : %s\n", sds.err)
	}

	m := new(module.Message)
	li := list.New()
	for sds.rows.Next() {
		err := sds.rows.StructScan(&m)
		if err != nil {
			log.Printf("SelectPushTargetStatus StructScan err : %s\n", sds.err)
		}
		li.PushBack(*m)
	}
	return li
}

// ExecMigrationBatchTarget ...
func (sds *SQLDataService) ExecMigrationBatchTarget(li *list.List) {
	defer func() {
		if r := recover(); r != nil {
			sds.tx.Rollback()
			log.Printf("[Recover] ExecMigrationBatch() : %s\n", r)
		}
	}()

	sds.db = module.DBconn(sds.DbURL)
	defer module.DBClose(sds.db)

	values := make([]string, 0, 0)
	for e := li.Front(); e != nil; e = e.Next() {
		values = append(values, e.Value.(module.Message).PushTargetSeq)
	}

	sds.db.MustExec(InsertPushTargetStatusLog + "(" + strings.Join(values, ",") + ")")
	sds.db.MustExec(InsertPushTargetLog + "(" + strings.Join(values, ",") + ")")

	// InnoDB 는 트랜잭션 처리한다.
	sds.tx = sds.db.MustBegin()
	sds.tx.MustExec(DeletePushTarget + "(" + strings.Join(values, ",") + ")")
	sds.tx.MustExec(DeletePushTargetStatus + "(" + strings.Join(values, ",") + ")")
	sds.tx.Commit()
}

// ExecMigrationBatchMessage ...
func (sds *SQLDataService) ExecMigrationBatchMessage() int64 {
	defer func() {
		module.DBClose(sds.db)
		if r := recover(); r != nil {
			log.Printf("[Recover] ExecMigrationBatchMessage() : %s\n", r)
		}
	}()
	sds.db = module.DBconn(sds.DbURL)
	result := sds.db.MustExec(DeletePushMessage)
	resultCnt, _ := result.RowsAffected()
	return resultCnt
}
