package utils

import (
	"errors"
	"strings"
)

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func HandleError(err error) string {
	if err != nil {
		return "Error: " + err.Error()
	}
	return ""
}

func JoinPaths(paths ...string) string {
	return strings.Join(paths, "/")
}

func SplitPath(path string) []string {
	return strings.Split(path, "/")
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func RemoveDuplicates(slice []string) []string {
	unique := make(map[string]struct{})
	result := []string{}

	for _, item := range slice {
		if _, exists := unique[item]; !exists {
			unique[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func ValidateFileName(name string) error {
	if IsEmpty(name) {
		return errors.New("file name cannot be empty")
	}
	if strings.Contains(name, "/") {
		return errors.New("file name cannot contain '/'")
	}
	return nil
}

