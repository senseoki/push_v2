package module

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// FileLog ...
type FileLog struct {
	Path string
}

// ExeFileLog ...
func (f *FileLog) ExeFileLog() {
	go func() {
	LOOP:
		t := time.Now()
		sYear := strconv.Itoa(t.Year())
		sMonth := strconv.Itoa(int(t.Month()))
		day := t.Day()
		sDay := strconv.Itoa(day)
		hour := strconv.Itoa(t.Hour())
		minute := strconv.Itoa(t.Minute())

		if len(sMonth) == 1 {
			sMonth = "0" + sMonth
		}
		if len(sDay) == 1 {
			sDay = "0" + sDay
		}
		if len(hour) == 1 {
			hour = "0" + hour
		}
		if len(minute) == 1 {
			minute = "0" + minute
		}

		filename := fmt.Sprintf("%s%s%s%s%s.log", sYear, sMonth, sDay, hour, minute)
		st := sYear + sMonth
		os.MkdirAll(f.Path+st+"/", 0777)
		logFile, err := os.OpenFile(f.Path+st+"/"+filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			log.Println("[로그파일 OPEN Error]")
		}
		log.SetOutput(io.MultiWriter(logFile, os.Stdout))
		for {
			if time.Now().Day() != day {
				logFile.Close()
				goto LOOP
			}
		}
	}()
}
