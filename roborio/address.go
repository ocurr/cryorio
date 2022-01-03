package roborio

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func GetAddresses(team int) []string {
	return []string{
		"roborio-" + strconv.FormatInt(int64(team), 10) + "-FRC.local",                                                            // roborio-TEAM-FRC.local
		"10." + strconv.FormatInt(int64(math.Floor(float64(team)/100)), 10) + "." + strconv.FormatInt(int64(team%100), 10) + ".2", // 10.TE.AM.2
		"172.22.11.2", // USB
	}
}

const (
	WpilibSettingsDir  = ".wpilib"
	WpilibSettingsFile = "wpilib_preferences.json"
)

func GetTeamNumber() (int, error) {
	wd, err := os.Getwd()
	if err != nil {
		return -1, err
	}

	settingsPath := reverseRecursiveFileSearch(wd, WpilibSettingsDir)
	if settingsPath == "" {
		return -1, fmt.Errorf("unable to find wpilib settings folder anywhere along: %s", wd)
	}

	settingsFile, err := os.ReadFile(filepath.Join(settingsPath, WpilibSettingsFile))
	if err != nil {
		return -1, fmt.Errorf("unable to open wpilib settings file: %w", err)
	}

	var settings map[string]interface{}
	err = json.Unmarshal(settingsFile, &settings)
	if err != nil {
		return -1, fmt.Errorf("unable to extract json from settings file: %w", err)
	}

	if team, ok := settings["teamNumber"]; ok {
		return int(team.(float64)), nil
	}

	return -1, fmt.Errorf("no team number present in repository")
}

// reverseRecursiveFileSearch searches the path specified by path for file.
// It expects path to have a basename that is a directory.
func reverseRecursiveFileSearch(path, file string) string {
	if len(path) <= 1 {
		return ""
	}

	entries, _ := os.ReadDir(path)
	for _, entry := range entries {
		if entry.Name() == file {
			return filepath.Join(path, entry.Name())
		}
	}

	return reverseRecursiveFileSearch(filepath.Dir(path), file)
}
