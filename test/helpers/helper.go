package helpers

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func ExecuteCliCommand(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()
	return strings.TrimSpace(buf.String()), err
}

func DeleteFile(fileName string) error {
	if err := os.Remove(fileName); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error cleaning up %s: %w", fileName, err)
	}
	return nil
}

func CopyFile(fileName, copyFileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	// Write data to dst
	return ioutil.WriteFile(copyFileName, data, 0755)
}

func isChar(char rune) bool {
	if 'a' <= char && char <= 'z' {
		return true
	}
	if 'A' <= char && char <= 'Z' {
		return true
	}
	return false
}

func KeepСhars(in string) string {
	buffer := make([]rune, 0, len(in))
	for _, char := range in {
		if isChar(char) {
			buffer = append(buffer, char)
		}
	}
	return string(buffer)
}

func CountNewlines(text string) int {
	count := 0
	for _, char := range text {
		if char == '\n' {
			count += 1
		}
	}
	return count
}
