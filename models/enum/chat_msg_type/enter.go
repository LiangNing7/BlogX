package chat_msg_type

type MsgType int8

const (
	TextMsgType     MsgType = 1
	ImageMsgType    MsgType = 2
	MarkdownMsgType MsgType = 3

	MsgReadType MsgType = 11
)
