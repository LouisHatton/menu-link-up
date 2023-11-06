package projects

import "time"

type Project struct {
	Id       string   `json:"id" firestore:"id"`
	Name     string   `json:"name" firestore:"name"`
	Slug     string   `json:"slug" firestore:"slug"`
	Metadata Metadata `json:"metadata" firestore:"metadata"`
	Users    []string `json:"users" firestore:"users"`
}

type Metadata struct {
	CreatedBy string    `json:"createdBy" firestore:"createdBy"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
}

type NewProject struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func Empty() *Project {
	return &Project{
		Users: []string{},
	}
}
