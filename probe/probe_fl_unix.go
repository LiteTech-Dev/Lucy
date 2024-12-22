//go:build unix || linux || darwin
// +build unix linux darwin

package probe

import (
	"errors"
	"os"
	"path"
	"syscall"
)

func checkServerFileLock() (locked bool, pid int) {
	lockPath := path.Join(
		getSavePath(),
		"session.lock",
	)
	file, err := os.OpenFile(lockPath, os.O_RDWR, 0666)
	defer file.Close()

	if err != nil {
		return false, 0
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if errors.Is(err, syscall.EWOULDBLOCK) {
		var fl syscall.Flock_t
		fl.Type = syscall.F_WRLCK
		fl.Whence = 0
		fl.Start = 0
		fl.Len = 0
		err = syscall.FcntlFlock(file.Fd(), syscall.F_GETLK, &fl)
		return true, int(fl.Pid)
	}
	syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

	return false, 0
}
