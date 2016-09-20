package main

import (
	"container/list"
	"flag"
	"log"
	"push_v2/batch/service"
	"push_v2/module"
	"sync"
	"time"
)

var (
	dbURL           string
	MQConnectionCnt int
	signalStatus    *module.SignalStatus
)

const (
	GoroutineCnt       = 20 // 1보다큰 짝수로 지정!
	sqlLimitPushTarget = 100000
	mqURL1             = "amqp://ezwel:ezwel@192.168.110.91:5672/push" //PushMQ01 (Batch)
	mqURL2             = "amqp://ezwel:ezwel@192.168.110.96:5672/push" //PushMQ03 (Batch)
	DEVDbURL           = "push:ezpush_0606@tcp(192.168.112.100:3306)/ez_push?charset=utf8"
	REALDbURL          = "push:ezpush_0606@tcp(192.168.111.23:3306)/ez_push?charset=utf8"
	pathLog            = "/pushlog/batch/"
	qName              = "push_queue"

	SENDSTATUSING      = "1002"
	SENDSTATUSCOMPLETE = "1003"
)

func main() {
	log.Println("[========== START PUSH BATCH V2 ==========]")
	// SqlDataService 생성하고 dbURL을 주입한다.
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
	// 1. 전송할 메세지를 조회한다.
	message := sqlDataService.GetMessage()
	if message.MsgSeq == "" {
		return
	}
	// 2. push_message status를 전송중(1002) 업데이트한다.
	log.Printf("[1] ServiceCd : %s  PushType : %s  MsgSeq : %s  발송중\n", message.ServiceCd, message.PushType, message.MsgSeq)
	sqlDataService.UpdateMessageSendStatus(message, SENDSTATUSING)

	// 3. MQ 전송 : MQ Connection
	mqConn, mqChSl := module.RunMQ(qName, mqURL1, MQConnectionCnt)
	mqConn2, mqChSl2 := module.RunMQ(qName, mqURL2, MQConnectionCnt)

	for {
		// 3-1. MQ 전송 : push_target 조회(전송할 사용자정보)
		listSlice := module.MakeSliceList(GoroutineCnt)
		sqlDataService.GetTargetUsers(listSlice, message)
		// 3-2. MQ 전송 : push_message status를 발송완료(1003) 업데이트한다. => 처음부터 다시 시작!
		if listSlice[GoroutineCnt-1].Len() == 0 {
			// 3-3. MQ close()
			module.CloseMQ(mqConn, mqChSl)
			module.CloseMQ(mqConn2, mqChSl2)
			// 3-4. UpdateMessageSendStatus & InsertPushMessageLog
			sqlDataService.UpdateMessageSendStatus(message, SENDSTATUSCOMPLETE)
			sqlDataService.InsertPushMessageLog(message)
			log.Printf("[2] ServiceCd : %s  PushType : %s  MsgSeq : %s  발송완료\n", message.ServiceCd, message.PushType, message.MsgSeq)
			break
		}

		var wg sync.WaitGroup
		var cntSl1, cntSl2 int
		dbConn := module.DBconn(dbURL)
		for i, ls := range listSlice {
			if ls.Len() > 0 {
				wg.Add(1)
				pushQueue := new(module.PushQueue)
				if i%2 == 0 {
					go func(ls *list.List, cntSl2 int, pushQueue *module.PushQueue) {
						confirmedSl := pushQueue.SendMQ(ls, mqChSl2[cntSl2])
						if len(confirmedSl) > 0 {
							service.InsertStatus(confirmedSl, dbConn)
						}
						wg.Done()
					}(ls, cntSl2, pushQueue)
					cntSl2++
				} else {
					go func(ls *list.List, cntSl1 int, pushQueue *module.PushQueue) {
						confirmedSl := pushQueue.SendMQ(ls, mqChSl[cntSl1])
						if len(confirmedSl) > 0 {
							service.InsertStatus(confirmedSl, dbConn)
						}
						wg.Done()
					}(ls, cntSl1, pushQueue)
					cntSl1++
				}
			}
		}
		wg.Wait()
		module.DBClose(dbConn)
	}
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
	// 명명줄 옵션: DEV, REAL 구분셋팅
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
