//go:build !windows

package roborio

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getTestRio(t *testing.T) *Roborio {
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
	return rio
}

func TestRoborio(t *testing.T) {
	t.Run("CreateAndSeeFiles", func(t *testing.T) {
		rio := getTestRio(t)
		out, err := rio.Touch("foo.txt")
		if err != nil {
			t.Fatal(err)
		}
		if len(out) != 0 {
			t.Fatal(out)
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
	t.Run("CreateAndDeleteFiles", func(t *testing.T) {
		rio := getTestRio(t)
		out, err := rio.Touch("bar.txt")
		if err != nil {
			t.Fatal(err)
		}
		if len(out) != 0 {
			t.Fatal(out)
		}
		out, err = rio.Remove("bar.txt")
		if err != nil {
			t.Fatal(err)
		}
		if len(out) != 0 {
			t.Fatal(out)
		}
		out, err = rio.ListDir()
		if err != nil {
			t.Fatal(err)
		}
		got := strings.TrimSpace(string(out))
		if got != "" {
			t.Fatalf("want: %q got: %s", "", got)
		}
	})
	t.Run("CopyFiles", func(t *testing.T) {
		rio := getTestRio(t)
		out, err := rio.Touch("FRCUserProgram")
		if err != nil {
			t.Fatal(err)
		}
		if len(out) != 0 {
			t.Fatal(out)
		}
		out, err = rio.Copy("FRCUserProgram", "match.backup")
		if err != nil {
			t.Fatal(err)
		}
		if len(out) != 0 {
			t.Fatal(out)
		}
		out, err = rio.ListDir()
		if err != nil {
			t.Fatal(err)
		}
		gotStr := strings.TrimSpace(string(out))
		got := strings.Split(gotStr, "\n")
		want := []string{"FRCUserProgram", "match.backup"}
		if !cmp.Equal(got, want) {
			t.Fatalf("want: %q got: %q", want, got)
		}
	})
}
