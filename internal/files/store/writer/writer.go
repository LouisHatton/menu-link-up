package writer

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/LouisHatton/menu-link-up/internal/files"
	fileStore "github.com/LouisHatton/menu-link-up/internal/files/store"
	"go.uber.org/zap"
)

var _ fileStore.Writer = (*Writer)(nil)

type Writer struct {
	l          *zap.Logger
	collection string
	db         *firestore.CollectionRef
}

func New(logger *zap.Logger, collection string, client *firestore.Client) (*Writer, error) {
	r := Writer{
		l:          logger,
		collection: collection,
	}
	r.db = client.Collection(collection)
	return &r, nil
}

func (r *Writer) Set(id string, file *files.File) error {
	logger := r.l.With(zap.String("fileId", id))

	logger.Debug("setting file doc")

	_, err := r.db.Doc(id).Set(context.TODO(), file)
	if err != nil {
		return fmt.Errorf("error getting file: %w", err)
	}
	logger.Debug("file doc set")

	return nil
}

func (w *Writer) Delete(id string) error {
	logger := w.l.With(zap.String("fileId", id))

	logger.Debug("deleting file doc")
	_, err := w.db.Doc(id).Delete(context.TODO())
	if err != nil {
		return fmt.Errorf("error deleting file doc: %w", err)
	}
	logger.Debug("file doc deleted")

	return nil
}
