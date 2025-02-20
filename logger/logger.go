/*
Copyright 2024 4rcadia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logger

import (
	"os"

	"github.com/emirpasic/gods/lists/singlylinkedlist"
)

var (
	debug     = false
	toConsole = false
)

func UseDebug() {
	debug = true
}

func UseConsoleOutput() {
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
