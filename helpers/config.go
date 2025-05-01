package helpers

import (
	"os"
	"strconv"
)

func GetFEUrl() string {
	feUrl := os.Getenv("FE_URL")
	if feUrl == "" {
		feUrl = "http://localhost:3000"
	}
	return feUrl
}

func GetMaxFileUploadSize() int64 {
	maxFileUploadSize, _ := strconv.ParseInt(os.Getenv("MAX_FILE_SIZE"), 10, 64)
	if maxFileUploadSize <= 0 {
		maxFileUploadSize = 5 * 1024 * 1024 // default 5 mb
	}
	return maxFileUploadSize
}
