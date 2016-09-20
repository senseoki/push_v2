package main

import (
	"flag"
	"log"
	"push_v2/migrationBatch/service"
	"push_v2/module"
	"time"
)

var (
	dbURL        string
	signalStatus *module.SignalStatus
)

const (
	DEVDbURL  = "push:ezpush_0606@tcp(192.168.112.100:3306)/ez_push?charset=utf8"
	REALDbURL = "push:ezpush_0606@tcp(192.168.111.23:3306)/ez_push?charset=utf8"
	pathLog   = "/pushlog/migrationBatch/"
)

func main() {
	log.Printf("[========== START PUSH MIGRATION BATCH V2 ==========]")
	done := make(chan bool)

	// Migration Batch goroutine : push_target
	go func() {
		sqlDataService := &service.SQLDataService{DbURL: dbURL}
		for {
			signalStatus.SignalChk()
			startTime := time.Now()
			RunTarget(sqlDataService)
			log.Printf("[RunTarget 최종 실행시간] %s\n\n", time.Since(startTime))
			time.Sleep(time.Millisecond * 500)
		}
	}()

	// Migration Batch goroutine : push_message
	go func() {
		sqlDataService := &service.SQLDataService{DbURL: dbURL}
		for {
			signalStatus.SignalChk()
			startTime := time.Now()
			RunMessage(sqlDataService)
			log.Printf("[RunMessage 최종 실행시간] %s\n\n", time.Since(startTime))
			time.Sleep(time.Minute * 10)
		}
	}()

	<-done
}

// RunTarget ...
func RunTarget(sqlDataService *service.SQLDataService) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recover] RunTarget() : %s\n", r)
		}
	}()
	li := sqlDataService.GetMigrationBatchTarget()
	if li.Len() == 0 {
		return
	}
	log.Printf("[0]	MIGRATION BATCH TARGET : 처리중...  %d건 \n", li.Len())
	sqlDataService.ExecMigrationBatchTarget(li)
	log.Printf("[1]	MIGRATION BATCH TARGET : 완료!\n")
}

// RunMessage ...
// 하루지난 push_message 삭제
// 이미 push_message_log에 이주 되어 있음(batch시 이주된다. 통계문제로 인해서 미리 이주하고 후 삭제)
func RunMessage(sqlDataService *service.SQLDataService) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recover] RunMessage() : %s\n", r)
		}
	}()
	resultCnt := sqlDataService.ExecMigrationBatchMessage()
	log.Printf("[0]	MIGRATION BATCH MESSAGE DELETE : %d 삭제 완료!\n", resultCnt)
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
}
