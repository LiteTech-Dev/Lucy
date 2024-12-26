package types

type LogLevel uint8
type LogItem struct {
	Level   LogLevel
	Content error
}
