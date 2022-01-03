package roborio

import (
	"strings"
	"testing"
)

func TestRoborio(t *testing.T) {
	dr := newDockerRio(t)
	defer dr.shutdown()
	rio, err := NewRoborio("test", "test", Addresses("localhost"), Team(973))
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

		rio.Remove("*")
	})
}
