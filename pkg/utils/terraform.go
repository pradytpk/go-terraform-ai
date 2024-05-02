package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// EndsWithTf Check the suffix ends with .tf
//
//	@param str
//	@return bool
func EndsWithTf(str string) bool {
	return strings.HasSuffix(str, ".tf")
}

// RandomName generates the random name
//
//	@return string
func RandomName() string {
	randomBytes := make([]byte, 5)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)
	return fmt.Sprintf("terraform-%s.tf", randomString)
}

// GetName Modifies the name
//
//	@param name
//	@return string
func GetName(name string) string {
	name = RemoveBlankLinesFromString(name)
	if EndsWithTf(name) {
		return name
	}
	return RandomName()
}

// TerraformPath Used to locate terraform executable files
//
//	@return string
//	@return error
func TerraformPath() (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "terraform")
	} else {
		cmd = exec.Command("which", "terraform")
	}
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running Init: %w", err)
	}
	return strings.TrimRight(string(output), "\n"), nil
}
