package roborio

import (
	"os"
	"os/exec"
	"strings"
)

type TestConn struct {
	dir string
}

func (c *TestConn) Close() error {
	return nil
}

func (c *TestConn) Exec(command string) ([]byte, error) {
	os.Chdir(c.dir)
	cmd := []string{}
	cmd = append(cmd, strings.Split(command, " ")...)
	return exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
}
