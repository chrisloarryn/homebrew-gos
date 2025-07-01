package common

import (
	"bufio"
	"os"
	"strings"
)

// AppendToFile appends content to a file
func AppendToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = file.WriteString(content)
	return err
}

// WriteToFile writes content to a file
func WriteToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// HasConfigContent checks if a file contains specific configuration content
func HasConfigContent(filename, content string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), content) {
			return true
		}
	}
	return false
}
