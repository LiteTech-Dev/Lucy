package logger

import (
	"fmt"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"os"
)

var debug = false

func SetDebug() {
	debug = true
}

var queue = singlylinkedlist.New()

func createLogFactory(level logLevel) func(content error) {
	return func(content error) {
		queue.Add(&logItem{Level: level, Content: content})
	}

}

var (
	Info = func(content any) {
		queue.Add(&logItem{Level: lInfo, Content: content})
	}
	Warning = createLogFactory(lWarning)
	Error   = createLogFactory(lError)
	Fatal   = func(content error) {
		defer os.Exit(1)
		createLogFactory(lFatal)(content)
		WriteAll()
	}
	Debug = func(content any) {
		if debug {
			queue.Add(&logItem{Level: lDebug, Content: content})
		}
	}
)

func WriteAll() {
	if queue.Empty() == false {
		_, _ = fmt.Fprintln(os.Stderr, "")
	}
	for queue.Empty() == false {
		pop()
	}
}
