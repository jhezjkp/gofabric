package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func ExistsInArray(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func FindRegex(regexPattern string, filename string) (string, error) {
	fileName := filename
	openaiPattern := regexPattern
	regex, err := regexp.Compile(openaiPattern)
	if err != nil {
		return "", err
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	matches := regex.FindStringSubmatch(string(file))
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", nil
}

func InsertIntoConfiguration(keyname string, apiKey string, createTables func()) (error) {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fileName:= filepath.Join(usrHome, ".config/fabric/.env")
	// check if file exist
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		createTables()
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	openaiPattern := fmt.Sprintf("%s=(.*)\n", keyname)
	regex, err := regexp.Compile(openaiPattern)
	if err != nil {
		return err
	}
	if regex.MatchString(string(file)) {
		result := regex.ReplaceAllString(string(file), fmt.Sprintf("%s=%s\n", keyname,apiKey))
		os.WriteFile(fileName, []byte(result), 0644)
	} else {
		result := string(file) + fmt.Sprintf("\n%s=%s\n", keyname,apiKey)
		os.WriteFile(fileName, []byte(result), 0644)
	}
	return nil
}