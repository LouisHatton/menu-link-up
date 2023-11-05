package store

import (
	"github.com/LouisHatton/menu-link-up/internal/db/query"
	"github.com/LouisHatton/menu-link-up/internal/files"
)

type Reader interface {
	Get(id string) (*files.File, error)
	GetByUrl(urlid string) (*files.File, error)
	Many(opts query.Options, wheres ...query.Where) (*[]files.File, error)
}

type Writer interface {
	Set(id string, file *files.File) error
	Delete(id string) error
}

type Manager struct {
	Reader
	Writer
}

func New(r Reader, w Writer) *Manager {
	return &Manager{r, w}
}
