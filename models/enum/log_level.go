package enum

type LogLevelType int8

const (
	LogInfoLevel LogLevelType = 1
	LogWarnLevel LogLevelType = 2
	LogErrLevel  LogLevelType = 3
)

func (l LogLevelType) String() string {
	switch l {
	case LogInfoLevel:
		return "info"
	case LogWarnLevel:
		return "warn"
	case LogErrLevel:
		return "error"
	}
	return ""
}
