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

func KeepСhars(in string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.TrimSpace(in), " ", "",
		), "\n", "",
	)
}
