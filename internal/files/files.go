package files

import (
	"fmt"
	"time"
)

type File struct {
	// Unique id which identifies the file (can't change)
	ID   string `json:"id"`
	Name string `json:"name"`

	UserId string `json:"userId"`

	// Unique Id which identifies the file - used for the url string (can change)
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	S3Region  string    `json:"s3Region"`
	S3Bucket  string    `json:"s3Bucket"`
	S3Key     string    `json:"s3Key"`
}

type NewFile struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func Empty() File {
	return File{}
}

func (f File) GetUrl() string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", f.S3Bucket, f.S3Region, f.S3Key)
}
