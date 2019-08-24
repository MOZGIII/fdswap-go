package fdswap_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/MOZGIII/fdswap"
)

func TestFdswap(t *testing.T) {
	// Prepare fd to replace.
	toReplace := os.Stdout

	// Prepare file to replace fd with.
	testFilePath := "/tmp/testfile.txt"
	replaceWith, err := os.OpenFile(testFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		t.Fatalf("error while preparing test file: %s", err)
	}
	defer replaceWith.Close() // nolint: errcheck

	// Swap fds.
	restoreHandle, err := fdswap.SwapFiles(toReplace, replaceWith)
	if err != nil {
		t.Fatalf("error while swapping fds: %s", err)
	}

	// Write to a fd that we replaced.
	fmt.Fprintln(toReplace, "my test output")

	// Restore the fd with the original.
	if err := restoreHandle.Restore(); err != nil {
		t.Fatalf("error while restoring fds: %s", err)
	}

	// Write to a fd that we replaced.
	fmt.Fprintln(toReplace, "this should be printed as usual")

	// Ensure that test file has the data as we expected.
	if err := replaceWith.Close(); err != nil {
		t.Fatalf("error while closing test file: %s", err)
	}

	testFileData, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("error while reading test file results: %s", err)
	}
	testFileDataString := string(testFileData)
	if testFileDataString != "my test output\n" {
		t.Fatalf("test file did not contain expected data: %q", testFileDataString)
	}
}
