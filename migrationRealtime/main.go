package main

import (
	"log"
	"push_v2/migrationRealtime/service"
	"push_v2/module"
	"time"
)

var (
	execMode     = "DEV" //DEV or REAL
	dbURL        string
	signalStatus *module.SignalStatus
)

const (
	DEVDbURL  = "push:ezpush_0606@tcp(192.168.112.100:3306)/ez_push?charset=utf8"
	REALDbURL = "push:ezpush_0606@tcp(192.168.112.23:3306)/ez_push?charset=utf8"
	pathLog   = "/pushlog/migrationRealtime/"
)

func main() {
	log.Printf("[========== START PUSH MIGRATION REALTIME ==========]")
	// SqlDataService 생성하고 dbURL을 주입한다.
	sqlDataService := &service.SQLDataService{DbURL: dbURL}
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
	li := sqlDataService.GetMigrationRealtimeMessage()
	if li.Len() == 0 {
		return
	}
	log.Printf("[0]	MIGRATION REALTIME : 처리중...  %d건 \n", li.Len())
	sqlDataService.ExecMigrationRealtime(li)
	log.Printf("[1]	MIGRATION REALTIME : 완료!\n")
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
	if execMode == "DEV" {
		dbURL = DEVDbURL
	} else {
		dbURL = REALDbURL
	}
}
