package files

import "time"

type File struct {
	// Unique id which identifies the file (can't change)
	Id string `json:"id" firestore:"id"`

	// Unique Id which identifies the file - used for the url string (can change)
	Slug      string   `json:"slug" firestore:"slug"`
	ProjectId string   `json:"projectId" firestore:"projectId"`
	Name      string   `json:"name" firestore:"name"`
	Metadata  Metadata `json:"metadata" firestore:"metadata"`
}

type Metadata struct {
	CreatedBy string    `json:"createdBy" firestore:"createdBy"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
}

type NewFile struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func Empty() File {
	return File{}
}
