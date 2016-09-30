package module

import (
	"container/list"
	"encoding/json"
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

type PushQueue struct {
	SendType      string `json:"SEND_TYPE`
	PushTargetSeq string `json:"PUSH_TARGET_SEQ"`
	ServiceCd     string `json:"SERVICE_CD"`
	PushType      string `json:"PUSH_TYPE"`
	SendMsg       string `json:"SEND_MSG"`
	PushToken     string `json:"PUSH_TOKEN"`
	OsCd          string `json:"OS_CD"`
	ImgTitle      string `json:"IMG_TITLE"`
	ImgFilePath   string `json:"IMG_FILE_PATH"`
	LinkUrl       string `json:"LINK_URL"`
	MsgType       string `json:"MSG_TYPE"`
	MsgSeq        string `json:"MSG_SEQ"`
	UserKey       string `json:"USER_KEY"`
	TestYn        string `json:"TEST_YN"`
}

type RepeatPushQueue struct {
	PushQueue
}

// RunMQ ...
func RunMQ(qName string, mqURL string, GoroutineCnt int) ([]*amqp.Connection, []*amqp.Channel) {
	var mqConnSl []*amqp.Connection
	var mqChSl []*amqp.Channel

	for {
		mqConnSl, mqChSl = ConnMQ(qName, mqURL, GoroutineCnt)
		if len(mqConnSl) > 0 {
			break
		} else {
			log.Printf("Not MQ Connection ... : %s", mqURL)
		}
	}
	return mqConnSl, mqChSl
}

// ConnMQ ...
func ConnMQ(qName string, mqURL string, GoroutineCnt int) ([]*amqp.Connection, []*amqp.Channel) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recover] MQconn : %s\n", r)
		}
	}()

	mqConnSl := make([]*amqp.Connection, GoroutineCnt)
	mqChSl := make([]*amqp.Channel, GoroutineCnt)

	for i := 0; i < GoroutineCnt; i++ {
		connMQ, err := amqp.Dial(mqURL)
		if err != nil {
			log.Printf("amqp.Dial() Error : %s\n", err)
		}
		ch, err := connMQ.Channel()
		if err != nil {
			log.Printf("connMQ.Channel() Error : %s\n", err)
		}

		_, err = ch.QueueDeclare(
			qName, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			true,  // no-wait
			nil,   // arguments
		)
		if err != nil {
			log.Printf("channel.QueueDeclare Error : %s\n", err)
		}
		mqChSl[i] = ch
		mqConnSl[i] = connMQ
	}
	return mqConnSl, mqChSl
}

// CloseMQ ...
func CloseMQ(connSl []*amqp.Connection, chSl []*amqp.Channel) {
	for _, ch := range chSl {
		ch.Close()
	}
	for _, conn := range connSl {
		conn.Close()
	}
	chSl = chSl[:0]
	connSl = connSl[:0]
}

// SendMQ (PushQueue)...
func (pushQueue *PushQueue) SendMQ(list *list.List, ch *amqp.Channel, sendType string) []string {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recover] PushQueue sendMQ() : %s\n", r)
		}
	}()
	confirmedSl := make([]string, 0, list.Len())
	for e := list.Front(); e != nil; e = e.Next() {
		pushQueue.SendType = sendType
		pushQueue.PushTargetSeq = e.Value.(Message).PushTargetSeq
		pushQueue.ServiceCd = e.Value.(Message).ServiceCd
		pushQueue.PushType = e.Value.(Message).PushType
		pushQueue.SendMsg = e.Value.(Message).SendMsg
		pushQueue.PushToken = e.Value.(Message).PushToken
		pushQueue.OsCd = e.Value.(Message).OsCd
		pushQueue.ImgTitle = e.Value.(Message).ImgTitle
		pushQueue.ImgFilePath = e.Value.(Message).ImgFilePath
		pushQueue.LinkUrl = e.Value.(Message).LinkUrl
		pushQueue.MsgType = e.Value.(Message).MsgType
		pushQueue.MsgSeq = e.Value.(Message).MsgSeq
		pushQueue.UserKey = e.Value.(Message).UserKey
		pushQueue.TestYn = e.Value.(Message).TestYn

		payload, _ := json.Marshal(pushQueue)
		err := ch.Publish(
			"",
			"push_queue", //routing key
			false,        //mandatory
			false,        //immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        payload,
			})
		if err != nil {
			log.Printf("channel.Publish Error : %s %s\n", e.Value.(Message).PushTargetSeq, err)
			break
		} else {
			confirmedSl = append(confirmedSl, "(\""+e.Value.(Message).PushTargetSeq+"\", \"1\")")
		}
	}
	return confirmedSl
}

// SendMQ (RepeatPushQueue)...
func (pushQueue *RepeatPushQueue) SendMQ(li *list.List, ch *amqp.Channel, sendType string) []string {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recover] RepeatPushQueue sendMQ() : %s\n", r)
		}
	}()
	confirmedSl := make([]string, 0, li.Len())
	for e := li.Front(); e != nil; e = e.Next() {
		pushQueue.SendType = sendType
		pushQueue.PushTargetSeq = e.Value.(Message).PushTargetSeq
		pushQueue.ServiceCd = e.Value.(Message).ServiceCd
		pushQueue.PushType = e.Value.(Message).PushType
		pushQueue.SendMsg = e.Value.(Message).SendMsg
		pushQueue.PushToken = e.Value.(Message).PushToken
		pushQueue.OsCd = e.Value.(Message).OsCd
		pushQueue.ImgTitle = e.Value.(Message).ImgTitle
		pushQueue.ImgFilePath = e.Value.(Message).ImgFilePath
		pushQueue.LinkUrl = e.Value.(Message).LinkUrl
		pushQueue.MsgType = e.Value.(Message).MsgType
		pushQueue.MsgSeq = e.Value.(Message).MsgSeq
		pushQueue.UserKey = e.Value.(Message).UserKey
		pushQueue.TestYn = e.Value.(Message).TestYn

		payload, _ := json.Marshal(pushQueue)

		statusNum, _ := strconv.Atoi(e.Value.(Message).MaxSendStatus)
		// 3 이면 다른 전송수단 선택(소스상에서)
		if statusNum < 3 {
			err := ch.Publish(
				"",
				"push_queue", //routing key
				false,        //mandatory
				false,        //immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        payload,
				})
			if err != nil {
				log.Printf("channel.Publish : %s %s\n", pushQueue.PushTargetSeq, err)
				break
			}
		} else {
			// To do 다른 전송수단
			// statusNum == 3
		}
		statusNum++
		confirmedSl = append(confirmedSl, "(\""+pushQueue.PushTargetSeq+"\", \""+strconv.Itoa(statusNum)+"\")")
	}
	return confirmedSl
}
