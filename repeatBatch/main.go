package main

import (
	"container/list"
	"flag"
	"log"
	batch "push_v2/batch/service"
	"push_v2/module"
	"push_v2/repeatBatch/service"
	"sync"
	"time"
)

var (
	dbURL           string
	MQConnectionCnt int
	signalStatus    *module.SignalStatus
)

const (
	SendType           = "1002" // Batch
	GoroutineCnt       = 20     // 1보다큰 짝수로 지정!
	sqlLimitPushTarget = 100000
	mqURL1             = "amqp://ezwel:ezwel@192.168.110.91:5672/push" //PushMQ01 (Batch)
	mqURL2             = "amqp://ezwel:ezwel@192.168.110.96:5672/push" //PushMQ03 (Batch)
	DEVDbURL           = "push:ezpush_0606@tcp(192.168.112.100:3306)/ez_push?charset=utf8&parseTime=true&loc=Local"
	REALDbURL          = "push:ezpush_0606@tcp(192.168.111.23:3306)/ez_push?charset=utf8&parseTime=true&loc=Local"
	pathLog            = "/pushlog/repeatBatch/"
	qName              = "push_queue"
)

func main() {
	log.Println("[========== START PUSH REPEAT-BATCH V2 ==========]")
	sqlDataService := &service.SQLDataService{DbURL: dbURL, SqlLimitPushTarget: sqlLimitPushTarget, GoroutineCnt: GoroutineCnt}
	for {
		signalStatus.SignalChk()
		startTime := time.Now()
		Run(sqlDataService)
		log.Printf("[REPEAT BATCH 최종 실행시간] %s\n\n", time.Since(startTime))
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
	sqlDataService.GetRepeatMessage(listSlice)
	if listSlice[GoroutineCnt-1].Len() == 0 {
		return
	}
	log.Printf("[0]\t Repeat 발송중\n")

	mqConn, mqChSl := module.RunMQ(qName, mqURL1, MQConnectionCnt)
	mqConn2, mqChSl2 := module.RunMQ(qName, mqURL2, MQConnectionCnt)

	var wg sync.WaitGroup
	var cntSl1, cntSl2 int
	dbConn := module.DBconn(dbURL)
	for i, ls := range listSlice {
		if ls.Len() > 0 {
			wg.Add(1)
			pushQueue := new(module.RepeatPushQueue)
			if i%2 == 0 {
				go func(ls *list.List, cntSl2 int, pushQueue *module.RepeatPushQueue) {
					confirmedSl := pushQueue.SendMQ(ls, mqChSl2[cntSl2], SendType)
					if len(confirmedSl) > 0 {
						batch.InsertStatus(confirmedSl, dbConn)
					}
					wg.Done()
				}(ls, cntSl2, pushQueue)
				cntSl2++
			} else {
				go func(ls *list.List, cntSl1 int, pushQueue *module.RepeatPushQueue) {
					confirmedSl := pushQueue.SendMQ(ls, mqChSl[cntSl1], SendType)
					batch.InsertStatus(confirmedSl, dbConn)
					wg.Done()
				}(ls, cntSl1, pushQueue)
				cntSl1++
			}
		}
	}
	wg.Wait()
	module.CloseMQ(mqConn, mqChSl)
	module.CloseMQ(mqConn2, mqChSl2)
	module.DBClose(dbConn)
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

	// MQ connection 갯수를 셋팅
	if GoroutineCnt < 2 {
		MQConnectionCnt = 1
	} else {
		MQConnectionCnt = GoroutineCnt / 2
	}
}
