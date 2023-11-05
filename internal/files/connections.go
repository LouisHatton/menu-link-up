package files

import "time"

type FileStatus string

const (
	Active    FileStatus = "active"
	Deploying FileStatus = "deploying"
	Paused    FileStatus = "paused"
	Unhealthy FileStatus = "unhealthy"
)

type File struct {
	// Unique id which identifies the file (can't change)
	Id string `json:"id" firestore:"id"`

	// Unique Id which identifies the file - used for the url string (can change)
	UrlId     string                 `json:"urlId" firestore:"urlId"`
	ProjectId string                 `json:"projectId" firestore:"projectId"`
	Name      string                 `json:"name" firestore:"name"`
	Metadata  Metadata               `json:"metadata" firestore:"metadata"`
	Tags      []string               `json:"tags" firestore:"tags"`
	Status    FileStatus             `json:"status" firestore:"status"`
	Schema    map[string]interface{} `json:"schema" firestore:"schema"`
}

type Metadata struct {
	CreatedBy string    `json:"createdBy" firestore:"createdBy"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
}

type NewFile struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func Empty() File {
	return File{
		Tags:   []string{},
		Schema: map[string]interface{}{},
	}
}
