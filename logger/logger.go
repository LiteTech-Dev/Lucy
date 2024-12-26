package logger

import (
	"fmt"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"lucy/output"
	"lucy/types"
	"os"
)

var messageQueue = singlylinkedlist.New()
var DisplayLevel = types.LogLevel(2)

func WriteAll() {
	for messageQueue.Empty() == false {
		popLogItem()
	}
	fmt.Println()
}

func CreateInfo(err error) {
	messageQueue.Add(&types.LogItem{Level: 0, Content: err})
}

func CreateWarning(err error) {
	messageQueue.Add(&types.LogItem{Level: 1, Content: err})
}

func CreateError(err error) {
	messageQueue.Add(&types.LogItem{Level: 2, Content: err})
}

func CreateFatal(err error) {
	defer os.Exit(1)
	messageQueue.Add(&types.LogItem{Level: 3, Content: err})
	WriteAll()
}

func popLogItem() {
	msg, _ := messageQueue.Get(0)
	if msg.(*types.LogItem).Level >= DisplayLevel {
		output.WriteLogItem(msg.(*types.LogItem))
	}
	messageQueue.Remove(0)
}
