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

// GetMessage ...
func (sc *SQLDataService) GetMessage() *module.Message {
	defer func() {
		module.DBClose(sc.db)
	}()
	sc.db = module.DBconn(sc.DbURL)
	m := new(module.Message)
	sc.err = sc.db.Get(m, Select_PushMessage)
	if sc.err != nil {
		//PushError(sc.err, "Select_PushMessage ... |")
	}
	return m
}

// UpdateMessageSendStatus ...
// send_status = 1002 전송중 업데이트.
// send_status = 1003 메세지는 하루가 지나면 migration 배치실행ㅇ후 push_message 테이블에서 삭제한다.
func (sc *SQLDataService) UpdateMessageSendStatus(m *module.Message, status string) {
	defer func() {
		module.DBClose(sc.db)
	}()

	sc.db = module.DBconn(sc.DbURL)
	if status == "1002" {
		_, sc.err = sc.db.Exec(Update_PushMessageSendStatus1002, status, m.ServiceCd, m.PushType, m.MsgSeq)
	} else {
		_, sc.err = sc.db.Exec(Update_PushMessageSendStatus1003, status, m.ServiceCd, m.PushType, m.MsgSeq)
	}

	if sc.err != nil {
		log.Printf("Update_PushMessageSendStatus : %s\n", sc.err)
	}
}

// InsertPushMessageLog ...
// (send_status = 1003 발송 완료 업데이트 후 log 테이블로 insert된다)
func (sc *SQLDataService) InsertPushMessageLog(m *module.Message) {
	defer func() {
		module.DBClose(sc.db)
		if r := recover(); r != nil {
			log.Printf("[Recover] InsertPushMessageLog : %s\n", r)
		}
	}()

	sc.db = module.DBconn(sc.DbURL)
	_, sc.err = sc.db.Exec(Insert_PushMessageLog, m.ServiceCd, m.PushType, m.MsgSeq)
	if sc.err != nil {
		log.Printf("Insert_PushMessageLog : %s\n", sc.err)
	}
}

// GetTargetUsers ...
func (sc *SQLDataService) GetTargetUsers(listSlice []*list.List, m *module.Message) {
	defer func() {
		sc.rows.Close()
		module.DBClose(sc.db)
	}()

	sc.db = module.DBconn(sc.DbURL)
	sc.rows, sc.err = sc.db.Queryx(Select_PushTarget, m.ServiceCd, m.PushType, m.MsgSeq, sc.SqlLimitPushTarget)
	if sc.err != nil {
		log.Printf("Select_PushTarget : %s\n", sc.err)
	}

	indexGoroutine := sc.GoroutineCnt
	for sc.rows.Next() {
		err := sc.rows.StructScan(&m)
		if err != nil {
			log.Printf("Select_PushTarget : %s\n", sc.err)
		}
		indexGoroutine--
		listSlice[indexGoroutine].PushBack(*m)
		if indexGoroutine == 0 {
			indexGoroutine = sc.GoroutineCnt
		}
	}
}

// InsertStatus ...
func InsertStatus(inVals []string, db *sqlx.DB) {
	_, err := db.Exec(Insert_PushTargetStatus + strings.Join(inVals, ","))
	if err != nil {
		log.Printf("InsertStatus : %s\n", err)
	}
}
