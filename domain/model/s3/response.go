package s3_model

import "time"

type UploadResponse struct {
	Key       string     `json:"key"`
	URL       string     `json:"url"`
	ExpiredAt *time.Time `json:"expiredAt"`
}
