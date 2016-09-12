package module

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type FileLog struct {
	Path string
}

// FileLog ...
func (f *FileLog) ExeFileLog() {
	go func() {
	LOOP:
		t := time.Now()
		sYear := strconv.Itoa(t.Year())
		sMonth := strconv.Itoa(int(t.Month()))
		day := t.Day()
		minute := t.Minute()

		filename := fmt.Sprintf("%s-%s-%d-%d-%d.log", sYear, sMonth, day, t.Hour(), minute)
		if len(sMonth) == 1 {
			sMonth = "0" + sMonth
		}
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
