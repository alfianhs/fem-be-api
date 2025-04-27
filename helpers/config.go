package helpers

import "os"

func GetFEUrl() string {
	feUrl := os.Getenv("FE_URL")
	if feUrl == "" {
		feUrl = "http://localhost:3000"
	}
	return feUrl
}
