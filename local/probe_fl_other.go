//go:build !windows && !unix && !linux && !darwin

package local

import "lucy/lucytypes"

func checkServerFileLock() *lucytypes.Activity {
	return nil
}
