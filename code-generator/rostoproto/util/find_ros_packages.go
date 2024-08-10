package util

import (
	"os/exec"
	"strings"
)

const (
	RosFindCommand = "rospack"
	RosFindArg     = "list"
)

func FindRosPackages() ([]string, error) {
	cmd := exec.Command(RosFindCommand, RosFindArg)
	output, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}
	packagePaths := []string{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		m := strings.IndexAny(line, " \t")
		if m < 0 {
			continue
		}
		line = strings.TrimSpace(line[m+1:])
		packagePaths = append(packagePaths, line)
	}
	return packagePaths, nil
}
