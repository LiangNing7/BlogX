package enum

type LogLevelType int8

const (
	LogInfoLevel LogLevelType = 1
	LogWarnLevel LogLevelType = 2
	LogErrLevel  LogLevelType = 3
)
