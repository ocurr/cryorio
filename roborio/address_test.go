package roborio

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetAddresses(t *testing.T) {
	cases := []struct {
		name   string
		input  int
		expect []string
	}{
		{
			"single digit team",
			1,
			[]string{
				"roborio-1-FRC.local",
				"10.0.1.2",
			},
		},
		{
			"two digit team",
			20,
			[]string{
				"roborio-20-FRC.local",
				"10.0.20.2",
			},
		},
		{
			"three digit team",
			973,
			[]string{
				"roborio-973-FRC.local",
				"10.9.73.2",
			},
		},
		{
			"four digit team",
			9973,
			[]string{
				"roborio-9973-FRC.local",
				"10.99.73.2",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			addrs := GetAddresses(c.input)
			// there are three addresses in the slice, however the last one never changes so we can ignore it.
			for i := 0; i < len(addrs)-1; i++ {
				if addrs[i] != c.expect[i] {
					t.Fatalf("incorrect rio address format: wanted: %s got: %s", c.expect[i], addrs[i])
				}
			}
		})
	}
}

func TestReverseRecursiveFileSearch(t *testing.T) {
	tempRoot := t.TempDir()
	fullPath := filepath.Join(tempRoot, "foo", "baz", "bar")
	wantPath := filepath.Join(tempRoot, "foo", "horse")

	err := os.MkdirAll(fullPath, 0777)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(wantPath, []byte{}, 0777)
	if err != nil {
		t.Fatal(err)
	}

	finalPath := reverseRecursiveFileSearch(fullPath, "horse")
	if finalPath != wantPath {
		t.Fatalf("error finding file: foundPath: %q actualPath: %q", finalPath, wantPath)
	}
}

func TestGetTeamNumber(t *testing.T) {
	tempRoot := t.TempDir()
	fullPath := filepath.Join(tempRoot, "src", "main", "cpp", "way", "too", "many", "dirs")
	dirPath := filepath.Join(tempRoot, WpilibSettingsDir)

	err := os.MkdirAll(fullPath, 0777)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dirPath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	wantPath := filepath.Join(dirPath, WpilibSettingsFile)
	contents, err := json.Marshal(
		struct {
			TeamNumber int `json:"teamNumber"`
		}{
			973,
		})
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(wantPath, contents, 0777)
	if err != nil {
		t.Fatal(err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir(fullPath)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(currentDir)
	team, err := GetTeamNumber()
	if err != nil {
		t.Fatal(err)
	}
	if team != 973 {
		t.Fatalf("wrong team number found: want: %d got: %d", 973, team)
	}
}
