package fileutil

import "os"

func ReadFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err == nil {
		return content, nil
	}
	content, err = os.ReadFile("../" + path)
	if err == nil {
		return content, nil
	}
	return nil, err
}
