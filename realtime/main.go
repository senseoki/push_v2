package main

import (
	"container/list"
	"flag"
	"log"
	"push_v2/module"
	"push_v2/realtime/service"
	"sync"
	"time"
)

var (
	dbURL        string
	signalStatus *module.SignalStatus
)

const (
	SendType           = "1001" // Real
	GoroutineCnt       = 60
	sqlLimitPushTarget = 10000
	mqURL1             = "amqp://ezwel:ezwel@192.168.110.155:5672/push" //PushRMQ01 (Realtime)
	DEVDbURL           = "push:ezpush_0606@tcp(192.168.112.100:3306)/ez_push?charset=utf8&parseTime=true&loc=Local"
	REALDbURL          = "push:ezpush_0606@tcp(192.168.111.23:3306)/ez_push?charset=utf8&parseTime=true&loc=Local"
	pathLog            = "/pushlog/realtime/"
	qName              = "push_queue"
)

func main() {
	log.Println("[========== START PUSH REALTIME V2 ==========]")
	// RealtimeSqlDataService 생성하고 dbURL을 주입한다.
	sqlDataService := &service.SQLDataService{DbURL: dbURL, SqlLimitPushTarget: sqlLimitPushTarget, GoroutineCnt: GoroutineCnt}
	for {
		signalStatus.SignalChk()
		startTime := time.Now()
		Run(sqlDataService)
		log.Printf("[REALTIME 최종 실행시간] %s\n\n", time.Since(startTime))
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
	sqlDataService.GetRealtimeMessage(listSlice)
	if listSlice[GoroutineCnt-1].Len() == 0 {
		return
	}

	log.Printf("[0]\t REALTIME 발송중\n")

	var wg sync.WaitGroup
	var cntSl int
	dbConn := module.DBconn(dbURL)
	// MQ Connection
	mqConn, mqChSl := module.RunMQ(qName, mqURL1, GoroutineCnt)
	for _, ls := range listSlice {
		if ls.Len() > 0 {
			wg.Add(1)
			pushQueue := new(module.PushQueue)
			go func(ls *list.List, cntSl int, pushQueue *module.PushQueue) {
				confirmedSl := pushQueue.SendMQ(ls, mqChSl[cntSl], SendType)
				if len(confirmedSl) > 0 {
					service.InsertRealtimeStatus(confirmedSl, dbConn)
				}
				wg.Done()
			}(ls, cntSl, pushQueue)
			cntSl++
		}
	}
	wg.Wait()
	module.CloseMQ(mqConn, mqChSl)
	module.DBClose(dbConn)
	log.Printf("[1]\t REALTIME 발송완료\n")
}

func init() {
	// #00. ExecSetting
	setExecSetting()
	// #01. Log File 을 설정하고 유지한다.
	fileLog := &module.FileLog{Path: pathLog}
	fileLog.ExeFileLog()
	// #02. 신호가 들어오면 프로그램 중지(현재 수행중인 프로세스 처리끝나면 중지)
	signalStatus = new(module.SignalStatus)
	signalStatus.InitSignal()
}

func setExecSetting() {
	// 명령줄 옵션: DEV, REAL 구분셋팅
	flagExecMode := flag.String("mode", "", "실행모드(DEV, REAL) 명령줄 옵션 없으면 기본 DEV 모드입니다.\n ex) -mode=DEV")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		*flagExecMode = "DEV"
	}

	if *flagExecMode == "DEV" {
		dbURL = DEVDbURL
	} else if *flagExecMode == "REAL" {
		dbURL = REALDbURL
	} else {
		*flagExecMode = "DEV"
		dbURL = DEVDbURL
	}
	log.Printf("프로그램 실행모드 : %s\n", *flagExecMode)
}
