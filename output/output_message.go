package output

import (
	"fmt"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"os"
)

var messageWriter = os.Stderr
var messageQueue = singlylinkedlist.New()
var DisplayLevel = messageLevel(2)

type messageLevel uint8
type message struct {
	Level   messageLevel
	Content error
}

func (s *message) PrintAndPop() {
	printMessage(s)
	index := messageQueue.IndexOf(*s)
	messageQueue.Remove(index)
}

func PrintMessagesAndExit(code int) {
	for messageQueue.Empty() == false {
		popMessage()
	}
	fmt.Println()
	os.Exit(code)
}

func CreateInfo(err error) {
	messageQueue.Add(&message{Level: 0, Content: err})
}

func CreateWarning(err error) {
	messageQueue.Add(&message{Level: 1, Content: err})
}

func CreateError(err error) {
	messageQueue.Add(&message{Level: 2, Content: err})
}

func CreateFatal(err error) {
	messageQueue.Add(&message{Level: 3, Content: err})
	PrintMessagesAndExit(1)
}

func popMessage() {
	msg, _ := messageQueue.Get(0)
	printMessage(msg.(*message))
	messageQueue.Remove(0)
}

func printMessage(message *message) {
	if message.Level >= DisplayLevel {
		fmt.Println()
		messageWriter.WriteString(formatMessage(message))
		fmt.Println()
	}
}

func formatMessage(message *message) string {
	switch message.Level {
	case 0:
		return cyan("[INFO] ") + message.Content.Error()
	case 1:
		return yellow("[WARNING] ") + message.Content.Error()
	case 2:
		return red("[ERROR] ") + message.Content.Error()
	case 3:
		return red("[FATAL] ") + message.Content.Error()
	}
	return "Wrong level"
}
