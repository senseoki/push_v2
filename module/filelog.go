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

		filename := fmt.Sprintf("%s.log", t.Format("200601021504"))
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
			time.Sleep(time.Millisecond * 1000)
		}
	}()
}
