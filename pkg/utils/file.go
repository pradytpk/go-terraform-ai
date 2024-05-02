package utils

import (
	"fmt"
	"log"
	"os"
)

// DirExists check if a directory exists at the specified path.
//
//	@param path
//	@return bool
func DirExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("Terraform is not initialized. Run `terraform init` first.")
	} else if err != nil {
		log.Fatalf("Failed to check if Terraform is initialized: %s\n", err.Error())
	}
	return true

}

// StoreFile writes the content to a file with the given name
//
//	@param name
//	@param contents
//	@return error
func StoreFile(name string, contents string) error {
	contents = RemoveBlankLinesFromString(contents)
	err := os.WriteFile(name, []byte(contents), 0o600)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil

}

// CurrenDir will find the current working directory
//
//	@return string
//	@return error
func CurrenDir() (string, error) {
	CurrenDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error CurrenDir:%w", err)
	}
	return CurrenDir, nil
}
