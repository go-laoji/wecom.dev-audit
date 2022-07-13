package archives

type StoreType string

const (
	DataBase StoreType = "1"
	Mq       StoreType = "2"
	Both     StoreType = "3"
)

type MsgType string

const (
	Text             MsgType = "text"
	Image            MsgType = "image"
	Revoke           MsgType = "revoke"
	Agree            MsgType = "agree"
	DisAgree         MsgType = "disagree"
	Voice            MsgType = "voice"
	Video            MsgType = "video"
	Card             MsgType = "card"
	Location         MsgType = "location"
	Emotion          MsgType = "emotion"
	File             MsgType = "file"
	Link             MsgType = "link"
	Weapp            MsgType = "weapp"
	Chatrecord       MsgType = "chatrecord"
	Collect          MsgType = "collect"
	Redpacket        MsgType = "redpacket"
	Meeting          MsgType = "meeting"
	Doc              MsgType = "doc"
	MarkDown         MsgType = "markdown"
	News             MsgType = "news"
	Calendar         MsgType = "calendar"
	Mixed            MsgType = "mixed"
	MeetingVoiceCall MsgType = "meeting_voice_call"
	VoipDocShare     MsgType = "voip_doc_share"
	Sphfeed          MsgType = "sphfeed"
)
