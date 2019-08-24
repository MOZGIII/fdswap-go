// +build !aix,!darwin,!dragonfly,!freebsd,!js,!wasm,!linux,!nacl,!netbsd,!openbsd,!solaris

package fdswap_test

import (
	"os"
	"syscall"
	"testing"

	"github.com/MOZGIII/fdswap-go"
)

func TestFdswap_other_produces_errors(t *testing.T) {
	_, err := fdswap.SwapFds(fdswap.FD(syscall.Stderr), fdswap.FD(syscall.Stdout))
	if err == nil {
		t.Errorf("expected an error on SwapFds, but didn't get any")
	}

	_, err = fdswap.SwapFiles(os.Stderr, os.Stdout)
	if err == nil {
		t.Errorf("expected an error on SwapFiles, but didn't get any")
	}
}
