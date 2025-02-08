package chat_msg

type TextMsg struct {
	Content string `json:"content"`
}
type ImageMsg struct {
	Href string `json:"href"`
}
type MarkdownMsg struct {
	Content string `json:"content"`
}
type ChatMsg struct {
	TextMsg     *TextMsg     `json:"textMsg,omitempty"`
	ImageMsg    *ImageMsg    `json:"imageMsg,omitempty"`
	MarkdownMsg *MarkdownMsg `json:"markdownMsg,omitempty"`
}
