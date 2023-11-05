package store

import (
	"github.com/LouisHatton/menu-link-up/internal/connections"
	"github.com/LouisHatton/menu-link-up/internal/db/query"
)

type Reader interface {
	Get(id string) (*connections.Connection, error)
	GetByUrl(urlid string) (*connections.Connection, error)
	Many(opts query.Options, wheres ...query.Where) (*[]connections.Connection, error)
}

type Writer interface {
	Set(id string, connection *connections.Connection) error
	Delete(id string) error
}

type Manager struct {
	Reader
	Writer
}

func New(r Reader, w Writer) *Manager {
	return &Manager{r, w}
}
