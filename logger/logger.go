package logger

import (
	"fmt"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"lucy/lucytypes"
	"lucy/output"
	"os"
)

var messageQueue = singlylinkedlist.New()
var DisplayLevel = lucytypes.LogLevel(2)

func WriteAll() {
	for messageQueue.Empty() == false {
		popLogItem()
	}
	fmt.Println()
}

func CreateInfo(err error) {
	messageQueue.Add(&lucytypes.LogItem{Level: 0, Content: err})
}

func CreateWarning(err error) {
	messageQueue.Add(&lucytypes.LogItem{Level: 1, Content: err})
}

func CreateError(err error) {
	messageQueue.Add(&lucytypes.LogItem{Level: 2, Content: err})
}

func CreateFatal(err error) {
	defer os.Exit(1)
	messageQueue.Add(&lucytypes.LogItem{Level: 3, Content: err})
	WriteAll()
}

func popLogItem() {
	msg, _ := messageQueue.Get(0)
	if msg.(*lucytypes.LogItem).Level >= DisplayLevel {
		output.WriteLogItem(msg.(*lucytypes.LogItem))
	}
	messageQueue.Remove(0)
}
