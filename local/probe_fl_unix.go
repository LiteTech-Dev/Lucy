//go:build unix || linux || darwin

package local

import (
	"errors"
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
		defer file.Close()

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
		syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

		return nil
	},
)
