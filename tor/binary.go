package tor

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"autonomia.digital/tonio/app/config"
)

const noPath = ""

// Initialize find a Tor binary that can be used by Tonio
func Initialize(configPath string) string {
	return findTorBinary(configPath)
}

func findTorBinary(configPath string) string {
	pathTorFound := checkInConfiguredPath(configPath)
	if pathTorFound != noPath {
		return pathTorFound
	}

	pathTorFound = checkInTonioDataDirectory()
	if pathTorFound != noPath {
		return pathTorFound
	}

	pathCWD, err := os.Getwd()
	if err == nil {
		pathTorFound = checkInLocalDirectory(pathCWD)
		if pathTorFound != noPath {
			return pathTorFound
		}

		pathTorFound = checkInExecutableDirectory(pathCWD)
		if pathTorFound != noPath {
			return pathTorFound
		}
	}

	pathTorFound = checkInCurrentWorkingDirectory()
	if pathTorFound != noPath {
		return pathTorFound
	}

	pathTorFound = checkInTonioBinary()
	if pathTorFound != noPath {
		return pathTorFound
	}

	pathTorFound = checkInHomeExecutableDirectory()
	if pathTorFound != noPath {
		return pathTorFound
	}

	pathTorFound = checkWithWhich()
	if pathTorFound != noPath {
		return pathTorFound
	}

	return noPath
}

func checkInConfiguredPath(configuredPath string) string {
	if isThereConfiguredTorBinary(configuredPath) {
		return configuredPath
	}
	return noPath
}

func checkInTonioDataDirectory() string {
	pathToFind := filepath.Join(config.XdgDataHome(), "tonio/tor")
	if isThereConfiguredTorBinary(pathToFind) {
		return pathToFind
	}
	return noPath
}

func checkInLocalDirectory(pathCWD string) string {
	pathToFind := filepath.Join(pathCWD, "/tor")
	if isThereConfiguredTorBinary(pathToFind) {
		return pathToFind
	}
	return noPath
}

func checkInExecutableDirectory(pathCWD string) string {
	pathToFind := filepath.Join(pathCWD, "/bin/tor")
	if isThereConfiguredTorBinary(pathToFind) {
		return pathToFind
	}
	return noPath
}

func checkInCurrentWorkingDirectory() string {
	pathToFind := filepath.Join(config.XdgDataHome(), "/tor")
	if isThereConfiguredTorBinary(pathToFind) {
		return pathToFind
	}
	return noPath
}

func checkInTonioBinary() string {
	pathToFind := filepath.Join(config.XdgDataHome(), "/bin/tonio/tor/tor")
	if isThereConfiguredTorBinary(pathToFind) {
		return pathToFind
	}
	return noPath
}

func checkInHomeExecutableDirectory() string {
	pathToFind := filepath.Join(config.XdgDataHome(), "/bin/tonio/tor")
	if isThereConfiguredTorBinary(pathToFind) {
		return pathToFind
	}
	return noPath
}

func checkWithWhich() string {
	outputWhich, err := executeCmd("which", []string{"tor"})
	if outputWhich == nil || err != nil {
		return noPath
	}

	pathToFind := strings.TrimSpace(string(outputWhich))
	if isThereConfiguredTorBinary(pathToFind) {
		return pathToFind
	}
	return noPath
}

func isThereConfiguredTorBinary(path string) bool {
	if path != noPath {
		return checkTorVersionCompatibility(path)
	}
	return false
}

func executeCmd(path string, args []string) ([]byte, error) {
	cmd := exec.Command(path, args...)
	output, err := cmd.Output()
	if output == nil || err != nil {
		return nil, errors.New("invalid command")
	}
	return output, nil
}

func checkTorVersionCompatibility(path string) bool {
	output, err := executeCmd(path, []string{"--version"})
	if output == nil || err != nil {
		return false
	}

	diff, err := compareVersions(extractVersionFrom(output), MinSupportedVersion)
	if err != nil {
		return false
	}

	return diff >= 0
}

func extractVersionFrom(s []byte) string {
	r := regexp.MustCompile(`(\d+\.)(\d+\.)(\d+\.)(\d)`)
	result := r.FindStringSubmatch(string(s))

	if len(result) == 0 {
		return ""
	}

	return result[0]
}
