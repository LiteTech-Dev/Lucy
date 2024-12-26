//go:build !windows && !unix && !linux && !darwin

package probe

import "lucy/lucytypes"

func checkServerFileLock() *lucytypes.Activity {
	return nil
}
