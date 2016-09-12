package module

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// SignalStatus ...
//var SignalStatus string
type SignalStatus struct {
	status string
}

// InitSignal register signals handler.
func (signalStatus *SignalStatus) InitSignal() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
			signalStatus.status = s.String()
			log.Println("SignalStatus :", signalStatus.status)
		case syscall.SIGHUP:
			// TODO reload
			//return
		default:
		}
	}()
}

// SignalChk ...
func (signalStatus *SignalStatus) SignalChk() {
	if signalStatus.status != "" {
		log.Printf("[프로그램 Signal 종료] : %s", signalStatus.status)
		os.Exit(0)
	}
}
