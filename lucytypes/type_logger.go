package lucytypes

type LogLevel uint8
type LogItem struct {
	Level   LogLevel
	Content error
}
