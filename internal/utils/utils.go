package utils

import (
	"bytes"
	"os/exec"
	"strings"
	"unicode/utf8"
)

func TrimFirstRune(s string) (rune, string) {
	r, i := utf8.DecodeRuneInString(s)
	return r, s[i:]
}

func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(out.String(), "\t", "   "), nil
}