package main

import (
	"flag"
	"log"
	"push_v2/migrationRealtime/service"
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
