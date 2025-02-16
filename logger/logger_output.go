package logger

import (
	"fmt"
	"log"
	"lucy/tools"
	"os"
)

func pop() {
	item, _ := queue.Get(0)
	writeItem(item.(*logItem))
	if debug && item.(*logItem).Level == lDebug {
		writeItem(item.(*logItem))
	}
	queue.Remove(0)
}

var color = map[logLevel]func(any) string{
	lInfo:    tools.Cyan,
	lWarning: tools.Yellow,
	lError:   tools.Red,
	lFatal:   tools.Red,
	lDebug:   tools.Green,
}

func (level logLevel) prefix(colored bool) string {
	return "[" + tools.Ternary(
		colored,
		color[level](level.String()),
		level.String(),
	) + "]"
}

func writeItem(message *logItem) {
	_, _ = fmt.Fprintln(os.Stderr, message.Level.prefix(true), message.Content)
	log.Println(message.Level.prefix(false), message.Content)
}
