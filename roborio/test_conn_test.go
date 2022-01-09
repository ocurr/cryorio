package roborio

import (
	"os"
	"os/exec"
	"strings"
)

type testConn struct {
	dir string
}

func (c *testConn) Close() error {
	return nil
}

func (c *testConn) Exec(command string) ([]byte, error) {
	os.Chdir(c.dir)
	cmd := []string{}
	cmd = append(cmd, strings.Split(command, " ")...)
	return exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
}
