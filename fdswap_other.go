// +build !aix,!darwin,!dragonfly,!freebsd,!js,!wasm,!linux,!nacl,!netbsd,!openbsd,!solaris

package fdswap

import (
	"errors"
	"os"
)

// FD is not supported on this platform.
type FD uintptr

var errNotImplemented = errors.New("not implemented")

// SwapFds is not supported on this platform.
func SwapFds(fdToReplace, fdToReplaceWith FD) (Restorer, error) {
	return nil, errNotImplemented
}

// SwapFiles is not supported on this platform.
func SwapFiles(fileToReplace, fileToReplaceWith *os.File) (Restorer, error) {
	return nil, errNotImplemented
}
