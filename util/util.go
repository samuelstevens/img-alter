package util

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"log"
	"os"
)

func HashFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()

	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	return base64.URLEncoding.EncodeToString(hash.Sum(nil)), nil
}
