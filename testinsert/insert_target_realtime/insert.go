package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	serviceCd = "1001" // 1001: ezwelfare
	pushType  = "1001" // 1001: ezadmin, 1002: 기념일
	msgSeq    = "1"
	osCd      = "20" //(00: 공통 ,10: iOS, 20:Android)
	sliceCnt  = 1
)

func main() {
	var msgSeq int
	for {
		msgSeq++
		fmt.Println("DB Insert Start...")
		startTime := time.Now()

		//db, err := sql.Open("mysql", "study:study@tcp(localhost:3306)/push?charset=utf8")
		db, err := sql.Open("mysql", "push:ezpush_0606@tcp(192.168.112.100:3306)/ez_push?charset=utf8")
		if err != nil {
			panic(err)
		}
		valueStrings := make([]string, 0, sliceCnt)
		for i := 0; i < sliceCnt; i++ {
			//valueStrings = append(valueStrings, "(\"message"+strconv.Itoa(i)+"\", 2001, \"TOKEN1459838497603\", 10, 32827980)")
			valueStrings = append(valueStrings, "('"+serviceCd+"', '"+pushType+"', \""+strconv.Itoa(msgSeq)+"\", \"1001\", \""+strconv.Itoa(msgSeq)+" ====== message2-1465446293315\", \"안드로이드 발송 테스트1465446293315\", \"/upload/2015/12/7/56id4sHd2v_20151207094927412000.png\", \"/index.jsp\", 32827980, \"A5UIZK8TECQ1HFUbyJap03EoJ2Kk5JCnKVk9S65YLIU=\",'"+osCd+"', \"e09bNqLHoEo:APA91bHpR4EbtOns_FHYlF962jiouLyJsjBw0QfFbVsZc7l80uJarzKfz5PzOQgyXQM5MK-j-TPaKjA-l5Xr2yKUjij8yVW7NPgl6ZDDPxEtCvy--WqjU9Hs3lOSa6r8fXO6LNi2uHXv\", NOW())")
		}
		_, err = db.Exec("INSERT INTO push_target_realtime (SERVICE_CD, PUSH_TYPE,MSG_SEQ,MSG_TYPE,SEND_MSG,IMG_TITLE,IMG_FILE_PATH,LINK_URL,USER_KEY,MOBILE,OS_CD,PUSH_TOKEN,REG_DT) VALUES " + strings.Join(valueStrings, ","))
		if err != nil {
			panic(err)
		}
		db.Close()

		//fmt.Printf("[INSERT] 최종 실행시간: %s\n", time.Since(startTime))
		fmt.Printf("[INSERT] 최종 실행시간: %s\n", time.Now().Sub(startTime))
		//time.Sleep(time.Millisecond * 100)
	}
}
