package logger

import (
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"os"
)

var debug = true
var toConsole = true

func SetDebug() {
	debug = true
}

func SetToConsole() {
	toConsole = true
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
	for queue.Empty() == false {
		pop()
	}
}
