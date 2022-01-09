//go:build !windows

package roborio

import (
	"strings"
	"testing"
)

func TestRoborio(t *testing.T) {
	rio, err := NewRoborio(func(user, pass, addr string) (Conn, error) {
		return &testConn{t.TempDir()}, nil
	}, "test", "test", Addresses("test"), Team(973), Port(8080))
	if err != nil {
		t.Fatal(err)
	}

	err = rio.Connect()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("CreateAndSeeFiles", func(t *testing.T) {
		out, err := rio.Touch("foo.txt")
		if err != nil {
			t.Fatal(err)
		}
		if len(out) != 0 {
			t.Fatal(err)
		}
		out, err = rio.ListDir()
		if err != nil {
			t.Fatal(err)
		}
		got := strings.TrimSpace(string(out))
		if got != "foo.txt" {
			t.Fatalf("want: foo.txt got: %s", got)
		}
	})
}
