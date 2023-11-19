package files

import (
	"time"
)

type File struct {
	// Unique id which identifies the file (can't change)
	ID   string `json:"id"`
	Name string `json:"name"`

	UserId string `json:"userId"`

	// Unique Id which identifies the file - used for the url string (can change)
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Slug      string    `json:"slug"`
	FileSize  int       `json:"fileSize"`
	S3Region  string    `json:"s3Region"`
	S3Bucket  string    `json:"s3Bucket"`
	S3Key     string    `json:"s3Key"`
}

type NewFile struct {
	Name     string `json:"name"`
	FileName string `json:"fileName"`
	Slug     string `json:"slug"`
	FileSize int    `json:"fileSize"`
}

type FileUpload struct {
	Url string `json:"url"`
}
