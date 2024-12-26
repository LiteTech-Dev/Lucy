//go:build windows
// +build windows

package probe

import (
	"golang.org/x/sys/windows"
	"lucy/lucytypes"
	"os"
	"path"
)

// This is AI generated code, please check it before use.
// I have no knowledge to Windows syscall.
var checkServerFileLock = memoize(
	func() *lucytypes.Activity {
		lockPath := path.Join(
			getSavePath(),
			"session.lock",
		)
		file, err := os.OpenFile(lockPath, os.O_RDWR, 0666)
		defer file.Close()

		if err != nil {
			return nil
		}

		err = windows.LockFileEx(
			windows.Handle(file.Fd()),
			windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
			0,
			1,
			0,
			&windows.Overlapped{},
		)
		if err != nil {
			var info windows.ByHandleFileInformation
			err = windows.GetFileInformationByHandle(
				windows.Handle(file.Fd()),
				&info,
			)
			if err == nil {
				return &lucytypes.Activity{
					Active: true,
					Pid:    int(info.VolumeSerialNumber),
				}
			}
		}

		windows.UnlockFileEx(
			windows.Handle(file.Fd()),
			0,
			1,
			0,
			&windows.Overlapped{},
		)

		return nil
	},
)
