package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

var (
	binName = "ginti-test"
)

func TestMain(m *testing.M) {
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binName, "../..")

	buf := &bytes.Buffer{}
	cmd.Stderr = buf

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build binary: ", err, buf.String())
		os.Exit(1)
	}

	result := m.Run()

	os.Remove(binName)
	os.Exit(result)
}
