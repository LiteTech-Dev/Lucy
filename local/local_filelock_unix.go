//go:build unix || darwin || linux

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

package local

import (
	"errors"
	"lucy/logger"
	"os"
	"path"
	"syscall"

	"lucy/lucytypes"
	"lucy/tools"
)

var checkServerFileLock = tools.Memoize(
	func() *lucytypes.Activity {
		if getSavePath() == "" {
			return nil
		}

		lockPath := path.Join(
			getSavePath(),
			"session.lock",
		)
		file, err := os.OpenFile(lockPath, os.O_RDWR, 0o666)
		defer tools.CloseReader(file, logger.Warning)

		if err != nil {
			return nil
		}

		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if errors.Is(err, syscall.EWOULDBLOCK) {
			var fl syscall.Flock_t
			fl.Type = syscall.F_WRLCK
			fl.Whence = 0
			fl.Start = 0
			fl.Len = 0
			err = syscall.FcntlFlock(file.Fd(), syscall.F_GETLK, &fl)
			return &lucytypes.Activity{
				Active: true,
				Pid:    int(fl.Pid),
			}
		}
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		if err != nil {
			return nil
		}

		return &lucytypes.Activity{
			Active: false,
			Pid:    0,
		}
	},
)
