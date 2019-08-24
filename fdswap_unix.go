// +build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

package fdswap

import (
	"os"
	"syscall"
)

// FD represents a file descriptor.
type FD = int

// SwapFds swaps the underlying target of `fdToReplace` with the one from
// `fdToReplaceWith`.
func SwapFds(fdToReplace, fdToReplaceWith FD) (*SwappedFdHandle, error) {
	syscall.ForkLock.RLock()
	defer syscall.ForkLock.RUnlock()

	// Dup the original Fd to be able to restore it.
	origFdCopy, err := syscall.Dup(fdToReplace)
	if err != nil {
		return nil, err
	}

	// Swap original fd with a new one.
	if err := syscall.Dup2(fdToReplaceWith, fdToReplace); err != nil {
		syscall.Close(origFdCopy) // nolint: errcheck, gosec
		return nil, err
	}

	return &SwappedFdHandle{
		fd:             fdToReplace,
		originalFdCopy: origFdCopy,
	}, nil
}

// SwapFiles swaps the underlying target of `fileToReplace` with the one from
// `fileToReplaceWith`.
func SwapFiles(fileToReplace, fileToReplaceWith *os.File) (*SwappedFdHandle, error) {
	return SwapFds(FD(fileToReplace.Fd()), FD(fileToReplaceWith.Fd()))
}

// SwappedFdHandle allows restoring swapped fd with original.
type SwappedFdHandle struct {
	fd             FD
	originalFdCopy FD
}

var _ Restorer = (*SwappedFdHandle)(nil)

// Restore underlying target of the fd to the original.
func (h *SwappedFdHandle) Restore() error {
	syscall.ForkLock.RLock()
	defer syscall.ForkLock.RUnlock()

	// Set the replaced fd back to it's original value.
	err := syscall.Dup2(h.originalFdCopy, h.fd)

	// Close the original fd copy to free system resources.
	// This call should not return an error cause we have another fd to the
	// underlying target.
	if closeErr := syscall.Close(h.originalFdCopy); closeErr != nil {
		panic(closeErr)
	}

	return err
}

// Fd returns Fd that was swapped.
func (h *SwappedFdHandle) Fd() FD {
	return h.fd
}

// OriginalFdCopy returns a Fd that is different from the original Fd, but has
// the same underlying target.
// Can be used to restore the Fd we swapped to it's original target.
func (h *SwappedFdHandle) OriginalFdCopy() FD {
	return h.originalFdCopy
}
