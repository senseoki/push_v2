package module

// Message ...
type Message struct {
	ServiceCd      string `db:"service_cd"`
	PushType       string `db:"push_type"`
	MsgSeq         string `db:"msg_seq"`
	MsgType        string `db:"msg_type"`
	SendMsg        string `db:"send_msg"`
	SendStatus     string `db:"send_status"`
	SendHopeDt     string `db:"send_hope_dt"`
	ImgTitle       string `db:"img_title"`
	ImgFilePath    string `db:"img_file_path"`
	LinkUrl        string `db:"link_url"`
	TotalCnt       string `db:"total_cnt"`
	IosSendCnt     string `db:"ios_send_cnt"`
	AndroidSendCnt string `db:"android_send_cnt"`
	RegDt          string `db:"reg_dt"`
	SendStartDt    string `db:"send_start_dt"`
	SendEndDt      string `db:"send_end_dt"`
	DelYn          string `db:"del_yn"`
	DelDt          string `db:"del_dt"`
	TestYn         string `db:"test_yn"`
	PushTargetSeq  string `db:"push_target_seq"`
	UserKey        string `db:"user_key"`
	Mobile         string `db:"mobile"`
	OsCd           string `db:"os_cd"`
	PushToken      string `db:"push_token"`
	TargetCount    int64
	MaxSendStatus  string `db:"max_send_status"`
	SchedulerWork  string `db:"scheduler_work"`
}
