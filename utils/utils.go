package utils

import (
	"fmt"
	"os"
	"regexp"
)

var (
	ErrGitURLFormat = fmt.Errorf("git url格式不对")
)

func ExtractGitInfo(gitURL string) (string, string, string, error) {
	re := regexp.MustCompile(`^git@(.*):(.*)/(.*).git$`)
	matches := re.FindStringSubmatch(gitURL)

	if len(matches) != 4 {
		return "", "", "", ErrGitURLFormat
	}

	return matches[1], matches[2], matches[3], nil
}

func AppendLineToFile(path, line string) {
	// Open the file in append mode
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error: open file %s failed\n", path)
		os.Exit(1)
	}
	defer file.Close()

	// Write the line to the file
	_, err = fmt.Fprintln(file, line)
	if err != nil {
		fmt.Printf("Error: write to file %s failed\n", path)
		os.Exit(1)
	}
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func IfDirExist(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return false
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
