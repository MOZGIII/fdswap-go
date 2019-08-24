package fdswap

import (
	"os"
	"syscall"
)

// SwapFds swaps the underlying target of `fdToReplace` with the one from
// `fdToReplaceWith`.
func SwapFds(fdToReplace, fdToReplaceWith int) (*SwappedFdHandle, error) {
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
	return SwapFds(int(fileToReplace.Fd()), int(fileToReplaceWith.Fd()))
}

// SwappedFdHandle allows restoring swapped fd with original.
type SwappedFdHandle struct {
	fd             int
	originalFdCopy int
}

// Restore underlying target of the fd to the original.
func (h *SwappedFdHandle) Restore() error {
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
func (h *SwappedFdHandle) Fd() int {
	return h.fd
}

// OriginalFdCopy returns a Fd that is different from the original Fd, but has
// the same underlying target.
// Can be used to restore the Fd we swapped to it's original target.
func (h *SwappedFdHandle) OriginalFdCopy() int {
	return h.originalFdCopy
}
