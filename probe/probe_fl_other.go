//go:build !windows && !unix && !linux && !darwin

package probe

import "lucy/types"

func checkServerFileLock() *types.Activity {
	return nil
}
