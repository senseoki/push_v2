package main

import (
	"container/list"
	"log"
	"push_v2/module"
	realtime "push_v2/realtime/service"
	"push_v2/repeatRealtime/service"
	"sync"
	"time"
)

var (
	execMode     = "DEV" //DEV or REAL
	dbURL        string
	signalStatus *module.SignalStatus
)

const (
	GoroutineCnt       = 20
	sqlLimitPushTarget = 100000
	mqURL1             = "amqp://ezwel:ezwel@192.168.110.155:5672/push" //PushRMQ01 (Realtime)
	DEVDbURL           = "push:ezpush_0606@tcp(192.168.112.100:3306)/ez_push?charset=utf8"
	REALDbURL          = "push:ezpush_0606@tcp(192.168.112.23:3306)/ez_push?charset=utf8"
	pathLog            = "/pushlog/repeatRealtime/"
	qName              = "push_queue"
)

func main() {
	log.Println("[========== START PUSH REPEAT-REALTIME V2 ==========]")
	// RealtimeSqlDataService 생성하고 dbURL을 주입한다.
	sqlDataService := &service.SQLDataService{DbURL: dbURL, SqlLimitPushTarget: sqlLimitPushTarget, GoroutineCnt: GoroutineCnt}
	for {
		signalStatus.SignalChk()
		startTime := time.Now()
		Run(sqlDataService)
		log.Printf("[최종 실행시간] %s\n\n", time.Since(startTime))
		time.Sleep(time.Millisecond * 500)
	}
}

// Run ...
func Run(sqlDataService *service.SQLDataService) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recover] Run() : %s\n", r)
		}
	}()
	listSlice := module.MakeSliceList(GoroutineCnt)
	// 전송할 메세지를 조회한다.
	sqlDataService.GetRepeatRealtimeMessage(listSlice)
	if listSlice[GoroutineCnt-1].Len() == 0 {
		return
	}

	log.Printf("[0]\t REPEAT-REALTIME 발송중\n")

	var wg sync.WaitGroup
	var cntSl int
	dbConn := module.DBconn(dbURL)
	// MQ Connection
	mqConn, mqChSl := module.RunMQ(qName, mqURL1, GoroutineCnt)
	for _, ls := range listSlice {
		if ls.Len() > 0 {
			wg.Add(1)
			pushQueue := new(module.RepeatPushQueue)
			go func(ls *list.List, cntSl int, pushQueue *module.RepeatPushQueue) {
				confirmedSl := pushQueue.SendMQ(ls, mqChSl[cntSl])
				if len(confirmedSl) > 0 {
					realtime.InsertRealtimeStatus(confirmedSl, dbConn)
				}
				wg.Done()
			}(ls, cntSl, pushQueue)
			cntSl++
		}
	}
	wg.Wait()
	module.CloseMQ(mqConn, mqChSl)
	module.DBClose(dbConn)
	log.Printf("[1]\t REPEAT-REALTIME 발송완료\n")
}

func init() {
	// #00. ExecSetting
	setExecSetting()
	// #01. Log File 을 설정하고 유지한다.
	fileLog := &module.FileLog{Path: pathLog}
	fileLog.ExeFileLog()
	// #02. 신호가 들어오면 프로그램 중지(현재 수행중인 프로세스 처리끝나면 중지)
	// 사용 : 배포나 유지보수용
	signalStatus = new(module.SignalStatus)
	signalStatus.InitSignal()
}

func setExecSetting() {
	if execMode == "DEV" {
		dbURL = DEVDbURL
	} else {
		dbURL = REALDbURL
	}
}
