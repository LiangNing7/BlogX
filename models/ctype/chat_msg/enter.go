package chat_msg

type TextMsg struct {
	Content string `json:"content"`
}
type ImageMsg struct {
	Src string `json:"src"`
}
type MarkdownMsg struct {
	Content string `json:"content"`
}
type MsgReadMsg struct {
	ReadChatID uint `json:"readChatID"` // 哪一条消息被读取了
}
type ChatMsg struct {
	TextMsg     *TextMsg     `json:"textMsg,omitempty"`
	ImageMsg    *ImageMsg    `json:"imageMsg,omitempty"`
	MarkdownMsg *MarkdownMsg `json:"markdownMsg,omitempty"`
	MsgReadMsg  *MsgReadMsg  `json:"msgReadMsg,omitempty"`
}
