package module

import (
	"log"
	"runtime"
	"time"
)

func ResolvePanic(str string, panicStatus *int) {
	if r := recover(); r != nil {
		if *panicStatus == 0 {
			*panicStatus = 1
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			slack := &Slack{Errmsg: string(buf[:n]), Info: str}
			slack.SendSlack()
		}
		log.Printf("%s : %s\n", str, r)
		time.Sleep(time.Millisecond * 1000)
	}
}
