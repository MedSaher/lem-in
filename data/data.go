package data

import (
	"bufio"
	"errors"
	"github.com/Douirat/lem-in/auth"
	"os"
)

// Exract data from the files as a single string to ease applying logic on it:
func ReadFile(fileName string) ([]string, error) {
	data := []string{}
	if !auth.IsValidFile(fileName) {
		return nil, errors.New("invalid file name")
	}
	file, err := os.Open("examples/" + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			data = append(data, line)
		}
	}
	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, errors.New("error reading file")
	}
	return data, nil
}
