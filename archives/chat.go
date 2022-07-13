package archives

type ChatData struct {
	Seq              uint64 `json:"seq"`
	MsgID            string `json:"msgid"`
	PublicKeyVer     uint32 `json:"publickey_ver"`
	EncryptRandomKey string `json:"encrypt_random_key"`
	EncryptChatMsg   string `json:"encrypt_chat_msg"`
}

type ChataDataResp struct {
	ErrCode  int        `json:"errcode"`
	ErrMsg   string     `json:"errmsg"`
	ChatData []ChatData `json:"chatdata"`
}

type MsgImage struct {
	Md5sum    string `json:"md5sum"`
	Sdkfileid string `json:"sdkfileid"`
	Filesize  uint32 `json:"filesize"`
}
type MsgRevoke struct {
	PreMsgid string `json:"pre_msgid"` // 标识撤回的原消息的msgid。String类型
}
type MsgAgree struct {
	Userid    string `json:"userid"`     // 同意/不同意协议者的userid，外部企业默认为external_userid。String类型
	AgreeTime int    `json:"agree_time"` // 同意/不同意协议的时间，utc时间，ms单位。
}
type MsgDisAgree struct {
	Userid       string `json:"userid"`
	DisagreeTime int    `json:"disagree_time"`
}
type MsgVoice struct {
	Md5sum     string `json:"md5sum"`
	VoiceSize  uint32 `json:"voice_size"`
	PlayLength uint32 `json:"play_length"`
	Sdkfileid  string `json:"sdkfileid"`
}
type MsgVideo struct {
	Md5Sum     string `json:"md5sum"`
	Filesize   int    `json:"filesize"`
	PlayLength int    `json:"play_length"`
	Sdkfileid  string `json:"sdkfileid"`
}
type MsgCard struct {
	Corpname string `json:"corpname"`
	Userid   string `json:"userid"`
}
type MsgLocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
	Title     string  `json:"title"`
	Zoom      int     `json:"zoom"`
}
type MsgEmotion struct {
	Type      int    `json:"type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Imagesize int    `json:"imagesize"`
	Md5Sum    string `json:"md5sum"`
	Sdkfileid string `json:"sdkfileid"`
}
type MsgFile struct {
	Md5Sum    string `json:"md5sum"`
	Filename  string `json:"filename"`
	Fileext   string `json:"fileext"`
	Filesize  int    `json:"filesize"`
	Sdkfileid string `json:"sdkfileid"`
}
type MsgLink struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	LinkURL     string `json:"link_url"`
	ImageURL    string `json:"image_url"`
}
type MsgWeapp struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Username    string `json:"username"`
	Displayname string `json:"displayname"`
}
type MsgChatrecord struct {
	Title string `json:"title"`
	Item  []struct {
		Type         string `json:"type"`
		Msgtime      int    `json:"msgtime"`
		Content      string `json:"content"`
		FromChatroom bool   `json:"from_chatroom"`
	} `json:"item"`
}
type MsgCollect struct {
	RoomName   string `json:"room_name"`
	Creator    string `json:"creator"`
	CreateTime string `json:"create_time"`
	Title      string `json:"title"`
	Details    []struct {
		ID   int    `json:"id"`
		Ques string `json:"ques"`
		Type string `json:"type"`
	} `json:"details"`
}
type MsgRedpacket struct {
	Type        int    `json:"type"`
	Wish        string `json:"wish"`
	Totalcnt    int    `json:"totalcnt"`
	Totalamount int    `json:"totalamount"`
}
type MsgMeeting struct {
	Topic       string `json:"topic"`
	Starttime   int    `json:"starttime"`
	Endtime     int    `json:"endtime"`
	Address     string `json:"address"`
	Remarks     string `json:"remarks"`
	Meetingtype int    `json:"meetingtype"`
	Meetingid   int    `json:"meetingid"`
	Status      int    `json:"status"`
}
type MsgDoc struct {
	Title      string `json:"title"`
	DocCreator string `json:"doc_creator"`
	LinkURL    string `json:"link_url"`
}
type MsgInfo struct {
	Content string     `json:"content,omitempty"` //　markdown时的消息内容
	Item    []struct { // 图文消息时才出现
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Picurl      string `json:"picurl"`
	} `json:"item,omitempty"`
}
type MsgCalendar struct {
	Title        string   `json:"title"`
	Creatorname  string   `json:"creatorname"`
	Attendeename []string `json:"attendeename"`
	Starttime    int      `json:"starttime"`
	Endtime      int      `json:"endtime"`
	Place        string   `json:"place"`
	Remarks      string   `json:"remarks"`
}
type MsgMixed struct {
	Item []struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"item"`
}
type MsgMeetingVoiceCall struct {
	Endtime      int    `json:"endtime"`
	Sdkfileid    string `json:"sdkfileid"`
	Demofiledata []struct {
		Filename     string `json:"filename"`
		Demooperator string `json:"demooperator"`
		Starttime    int    `json:"starttime"`
		Endtime      int    `json:"endtime"`
	} `json:"demofiledata"`
	Sharescreendata []struct {
		Share     string `json:"share"`
		Starttime int    `json:"starttime"`
		Endtime   int    `json:"endtime"`
	} `json:"sharescreendata"`
}
type MsgVoipDocShare struct {
	Filename  string `json:"filename"`
	Md5Sum    string `json:"md5sum"`
	Filesize  int    `json:"filesize"`
	Sdkfileid string `json:"sdkfileid"`
}
type MsgSphfeed struct {
	FeedType int    `json:"feed_type"`
	SphName  string `json:"sph_name"`
	FeedDesc string `json:"feed_desc"`
}
type ChataMsg struct {
	MsgId   string   `json:"msgid"`
	Action  string   `json:"action"`
	From    string   `json:"from"`
	ToList  []string `json:"tolist"`
	RoomId  string   `json:"roomid"`
	MsgTime int64    `json:"msgtime"`
	MsgType string   `json:"msgtype"`
	// 文本消息
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	//　图片
	Image MsgImage `json:"image"`

	//　撤回消息
	Revoke MsgRevoke `json:"revoke"`
	//　同意会话内容
	Agree MsgAgree `json:"agree"`
	//　不同意会话内容
	DisAgree MsgDisAgree `json:"disagree"`
	//　语音
	Voice MsgVoice `json:"voice"`
	//　视频
	Video MsgVideo `json:"video"`
	//　名片
	Card MsgCard `json:"card"`
	//　位置
	Location MsgLocation `json:"location"`
	//　表情
	Emotion MsgEmotion `json:"emotion"`
	//　文件
	File MsgFile `json:"file"`
	//　连接
	Link MsgLink `json:"link"`
	//　小程序消息
	Weapp MsgWeapp `json:"weapp"`
	//　会话记录消息
	Chatrecord MsgChatrecord `json:"chatrecord"`
	//	TODO 待办消息　投票消息
	//　填表消息
	Collect MsgCollect `json:"collect"`
	//	红包消息 & 互通红包消息
	Redpacket MsgRedpacket `json:"redpacket"`
	//	会议邀请信息
	Meeting MsgMeeting `json:"meeting"`
	//	TODO　切换企业日志　参考如下内容，暂不处理
	//	{"msgid":"125289002219525886280","action":"switch","time":1554119421840,"user":"XuJinSheng"}
	//　在线文档消息
	Doc MsgDoc `json:"doc"`
	//　MarkDown格式消息＆图文消息
	Info MsgInfo `json:"info"`
	//　日程消息
	Calendar MsgCalendar `json:"calendar"`
	Mixed    MsgMixed    `json:"mixed"`
	// 音频存档消息
	MeetingVoiceCall MsgMeetingVoiceCall `json:"meeting_voice_call"`
	// 音频共享文档消息
	VoipDocShare MsgVoipDocShare `json:"voip_doc_share,omitempty"`
	// 视频号消息
	Sphfeed MsgSphfeed `json:"sphfeed"`
}
