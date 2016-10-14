package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

const (
	serviceCd = "1001" // 1001: ezwelfare
	pushType  = "1001" // 1001: ezadmin, 1002: 기념일
	msgSeq    = "2"
	osCd      = "20" //(00: 공통 ,10: iOS, 20:Android)
	sliceCnt  = 1000000
)

func main() {
	//db, err := sqlx.Connect("mysql", "study:study@tcp(localhost:3306)/push?charset=utf8")
	//db, err := sqlx.Connect("mysql", "push:ezpush_0606@tcp(192.168.112.100:3306)/ez_push?charset=utf8&parseTime=true&loc=Local") // DEV
	db, err := sqlx.Connect("mysql", "push:ezpush_0606@tcp(192.168.111.23:3306)/ez_push?charset=utf8&parseTime=true&loc=Local") // REAL
	db.SetMaxOpenConns(100)
	tx := db.MustBegin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[Recover] main() : %s\n", r)
		}
		tx.Commit()
		db.Close()
	}()
	fmt.Println("DB Insert Target Batch Start...")
	startTime := time.Now()

	valueStrings := make([]string, 0, sliceCnt)
	for i := 0; i < sliceCnt; i++ {
		valueStrings = append(valueStrings, "('"+serviceCd+"', '"+pushType+"', "+msgSeq+", "+strconv.Itoa(i)+", \"A5UIZK8TECQ1HFUbyJap03EoJ2Kk5JCnKVk9S65YLIU=\", '"+osCd+"', \"cgTJ-7he25I:APA91bGDupKHBKHzM2PLdBbuGQmaQL6MskCIA-D0LRi39PLbs8Dxa-ou0KZuTrOwySfkibCAIdAcoc2DdtHpCkYnAC6TxuqGrK8I8GHaUCsFaopnDFH0NpW57UrPholvqK_kAPBVOBdG\", NOW())")
	}

	if err != nil {
		panic(err)
	}

	// SERVICE_CD=1001(ezwelfare), PUSH_TYPE=1001(ezadmin) 1002(기념일), MSG_TYPE=1002(이미지형)

	tx.MustExec(`INSERT INTO push_message (SERVICE_CD ,PUSH_TYPE, MSG_SEQ, MSG_TYPE, SEND_MSG, SEND_STATUS, SEND_HOPE_DT, IMG_TITLE, IMG_FILE_PATH, LINK_URL,
	                                            TOTAL_CNT, IOS_SEND_CNT, ANDROID_SEND_CNT, REG_DT, SEND_START_DT, SEND_END_DT, DEL_YN, DEL_DT, TEST_YN)
					  VALUES ('` + serviceCd + `', '` + pushType + `', '` + msgSeq + `', '1002', 'Batch Go 에서 넣은 Message입니다.', '1001', '20160822140000', '이미지', 'https:\/\/img.ezwelfare.net\/welfare\/upload\/2016\/6\/23\/715XZdw4FN_20160623145934105000.jpg', '',
					          0, 0, 0, NOW(), '', '', '', '', 'Y')`)
	tx.MustExec("INSERT INTO push_target (SERVICE_CD ,PUSH_TYPE ,MSG_SEQ, USER_KEY, MOBILE, OS_CD, PUSH_TOKEN, REG_DT) VALUES " + strings.Join(valueStrings, ","))

	fmt.Printf("[INSERT] 최종 실행시간: %s  , serviceCd: %s, pushType: %s, msgSeq: %s, \n", time.Since(startTime), serviceCd, pushType, msgSeq)
	//fmt.Printf("[INSERT] 최종 실행시간: %s\n", time.Now().Sub(startTime))
}
