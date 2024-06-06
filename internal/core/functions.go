package core

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
)

var imageTypes = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpeg",
	"image/webp": ".webp",
}

func GenerateRandomString(length int) (string, error) {
	randBytes := make([]byte, length)

	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(randBytes)
	if len(randomString) > length {
		randomString = randomString[:length]
	}

	return randomString, nil
}

func CreateFile(bytes []byte) (string, error) {
	filename, err := GenerateRandomString(64)
	if err != nil {
		return "", err
	}

	mime := http.DetectContentType(bytes)
	ext, ok := imageTypes[mime]
	if ok {
		filename += ext
	} else {
		return "", ErrUnsupportedImageFormat
	}

	file, err := os.Create("/uploads/" + filename)
	if err != nil {
		return "", err
	}

	_, err = file.Write(bytes)
	if err != nil {
		return "", err
	}

	return filename, err
}
