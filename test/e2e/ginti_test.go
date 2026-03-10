package e2e

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func getCommand(args ...string) (*exec.Cmd, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, binName)

	cmd := exec.Command(path, args...)
	return cmd, nil
}

func normalizeOutputLine(output string) string {
	output = strings.ReplaceAll(output, "\t", " ")
	fields := strings.Fields(output)
	return strings.Join(fields, " ")
}

func TestStdin(t *testing.T) {
	cmd, err := getCommand()
	if err != nil {
		t.Fatal("couldn't get working directory:", err)
	}

	output := &bytes.Buffer{}

	cmd.Stdin = strings.NewReader("one two three\n")
	cmd.Stdout = output

	if err := cmd.Run(); err != nil {
		t.Fatal("failed to run command")
	}

	wants := "1 3 14"
	got := normalizeOutputLine(output.String())

	if wants != got {
		t.Log("stdout is not correct: wanted:", wants, "got:", got)
		t.Fail()
	}
}

func TestSingleFile(t *testing.T) {
	file, err := os.CreateTemp("", "ginti-test-*")
	if err != nil {
		t.Fatal("couldn't create temp file")
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString("foo bar baz\nbaz bar foo\none two three\n")
	if err != nil {
		t.Fatal("could not write to temp file")
	}

	if err := file.Close(); err != nil {
		t.Fatal("could not close file")
	}

	cmd, err := getCommand(file.Name())
	if err != nil {
		t.Fatal("couldn't get working directory:", err)
	}

	output := &bytes.Buffer{}
	cmd.Stdout = output

	if err := cmd.Run(); err != nil {
		t.Fatal("failed to run command:", err)
	}

	wants := fmt.Sprintf("3 9 38 %s", file.Name())
	got := normalizeOutputLine(output.String())

	if got != wants {
		t.Log("stdout not expected: wants:", wants, "got:", got)
		t.Fail()
	}
}

func TestNoExist(t *testing.T) {
	cmd, err := getCommand("noexist.txt")
	if err != nil {
		t.Fatal("could not create command:", err)
	}

	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err = cmd.Run()
	if err == nil {
		t.Log("command succeeded when should have failed")
		t.Fail()
	}

	if err.Error() != "exit status 1" {
		t.Log("expected error of exit status 1. got:", err.Error())
		t.Fail()
	}

	stderrStr := stderr.String()
	if !strings.HasPrefix(stderrStr, "counter: open noexist.txt:") {
		t.Log("stderr not match prefix: got:", stderrStr)
		t.Fail()
	}

	if !strings.Contains(stderrStr, "no such file or directory") {
		t.Log("stderr not match OS error: got:", stderrStr)
		t.Fail()
	}

	if stdout.String() != "" {
		t.Log("stdout not match: wants empty got:", stdout.String())
		t.Fail()
	}
}

func createFile(content string) (*os.File, error) {
	file, err := os.CreateTemp("", "ginti-test-*")
	if err != nil {
		return nil, fmt.Errorf("could not create temp file")
	}
	_, err = file.WriteString(content)
	if err != nil {
		return nil, fmt.Errorf("could not write to temp file")
	}

	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("could not close file")
	}

	return file, nil
}

func TestMultipleFiles(t *testing.T) {
	fileA, err := createFile("one two three\n")
	if err != nil {
		t.Fatal("could not create file one:", err)
	}
	defer os.Remove(fileA.Name())

	fileB, err := createFile("foo bar\nbaz foo\n")
	if err != nil {
		t.Fatal("could not create file two:", err)
	}
	defer os.Remove(fileB.Name())

	fileC, err := createFile("")
	if err != nil {
		t.Fatal("could not create file three:", err)
	}
	defer os.Remove(fileC.Name())

	cmd, err := getCommand(fileA.Name(), fileB.Name(), fileC.Name())
	if err != nil {
		t.Fatal("couldn't get working directory:", err)
	}

	output := &bytes.Buffer{}
	cmd.Stdout = output

	if err := cmd.Run(); err != nil {
		t.Fatal("failed to run command:", err)
	}

	expected := map[string]string{
		fileA.Name(): fmt.Sprintf("1 3 14 %s", fileA.Name()),
		fileB.Name(): fmt.Sprintf("2 4 16 %s", fileB.Name()),
		fileC.Name(): fmt.Sprintf("0 0 0 %s", fileC.Name()),
		"total":     "3 7 30 total",
	}

	scanner := bufio.NewScanner(output)

	for scanner.Scan() {
		line := normalizeOutputLine(scanner.Text())
		if line == "" {
			t.Log("empty line detected")
			t.Fail()
			continue
		}

		fields := strings.Fields(line)
		filename := fields[len(fields)-1]

		expects, ok := expected[filename]
		if !ok {
			t.Log("no expectation for line:", line)
			t.Fail()
			continue
		}

		if expects != line {
			t.Log("line doesn't match, wants:", expects, "got:", line)
			t.Fail()
			continue
		}

		delete(expected, filename)
	}

	if err := scanner.Err(); err != nil {
		t.Log("scanner error:", err)
		t.Fail()
	}

	if len(expected) != 0 {
		for _, missing := range expected {
			t.Log("missing expected line:", missing)
			t.Fail()
		}
	}
}

func TestFlags(t *testing.T) {
	file, err := createFile("one two three four five\none two three\n")
	if err != nil {
		t.Fatal("could not create file")
	}
	defer os.Remove(file.Name())

	testCases := []struct {
		name  string
		flags []string
		want  string
	}{
		{
			name:  "line flag",
			flags: []string{"-l"},
			want:  fmt.Sprintf("2 %s", file.Name()),
		},
		{
			name:  "bytes flag",
			flags: []string{"-c"},
			want:  fmt.Sprintf("38 %s", file.Name()),
		},
		{
			name:  "words flag",
			flags: []string{"-w"},
			want:  fmt.Sprintf("8 %s", file.Name()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputs := append(tc.flags, file.Name())

			cmd, err := getCommand(inputs...)
			if err != nil {
				t.Log("failed to get command:", err)
				t.Fail()
			}

			stdout := &bytes.Buffer{}
			cmd.Stdout = stdout

			if err := cmd.Run(); err != nil {
				t.Log("failed to run command:", err)
				t.Fail()
			}

			output := normalizeOutputLine(stdout.String())

			if output != tc.want {
				t.Log("output did not match:", output, "wants:", tc.want)
				t.Fail()
			}
		})
	}
}
