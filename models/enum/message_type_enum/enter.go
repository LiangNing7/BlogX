package message_type_enum

type Type int8

const (
	CommentType          Type = 1
	ApplyType            Type = 2
	DiggArticleType      Type = 3
	UnDiggArticleType    Type = 4
	DiggCommentType      Type = 5
	UnDiggCommentType    Type = 6
	CollectArticleType   Type = 7
	UnCollectArticleType Type = 8
	SystemType           Type = 9
)
