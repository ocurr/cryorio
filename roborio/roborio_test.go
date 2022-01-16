//go:build !windows

package roborio

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getTestRio(t *testing.T) *Roborio {
	rio, err := NewRoborio(func(user, pass, addr string) (Conn, error) {
		return &TestConn{t.TempDir()}, nil
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

func assertCmdSuccess(t *testing.T, err error) {
	assertCmdSuccessAndNoOutput(t, []byte{}, err)
}

func assertCmdSuccessAndNoOutput(t *testing.T, out []byte, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 {
		t.Fatal(out)
	}
}

func TestRoborio(t *testing.T) {
	t.Run("CreateAndSeeFiles", func(t *testing.T) {
		rio := getTestRio(t)
		out, err := rio.Touch("foo.txt")
		assertCmdSuccessAndNoOutput(t, out, err)
		out, err = rio.ListDir()
		assertCmdSuccess(t, err)
		got := strings.TrimSpace(string(out))
		if got != "foo.txt" {
			t.Fatalf("want: foo.txt got: %s", got)
		}
	})
	t.Run("CreateAndDeleteFiles", func(t *testing.T) {
		rio := getTestRio(t)
		out, err := rio.Touch("bar.txt")
		assertCmdSuccessAndNoOutput(t, out, err)
		out, err = rio.Remove("bar.txt")
		assertCmdSuccessAndNoOutput(t, out, err)
		out, err = rio.ListDir()
		assertCmdSuccess(t, err)
		got := strings.TrimSpace(string(out))
		if got != "" {
			t.Fatalf("want: %q got: %s", "", got)
		}
	})
	t.Run("CopyFiles", func(t *testing.T) {
		rio := getTestRio(t)
		out, err := rio.Touch("FRCUserProgram")
		assertCmdSuccessAndNoOutput(t, out, err)
		out, err = rio.Copy("FRCUserProgram", "match.backup")
		assertCmdSuccessAndNoOutput(t, out, err)
		out, err = rio.ListDir()
		assertCmdSuccess(t, err)
		gotStr := strings.TrimSpace(string(out))
		got := strings.Split(gotStr, "\n")
		want := []string{"FRCUserProgram", "match.backup"}
		if !cmp.Equal(got, want) {
			t.Fatalf("want: %q got: %q", want, got)
		}
	})
	t.Run("BackupFile", func(t *testing.T) {
		rio := getTestRio(t)
		out, err := rio.Touch("FRCUserProgram")
		assertCmdSuccessAndNoOutput(t, out, err)
		err = rio.BackupFile("FRCUserProgram", "match")
		if err != nil {
			t.Fatal(err)
		}
		out, err = rio.ListDir()
		assertCmdSuccess(t, err)
		gotStr := strings.TrimSpace(string(out))
		got := strings.Split(gotStr, "\n")
		want := []string{"FRCUserProgram", "match.backup"}
		if !cmp.Equal(got, want) {
			t.Fatalf("want: %q got: %q", want, got)
		}
	})
}
