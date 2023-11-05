package reader

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	dbFirestore "github.com/LouisHatton/menu-link-up/internal/db/firestore"
	"github.com/LouisHatton/menu-link-up/internal/db/query"
	"github.com/LouisHatton/menu-link-up/internal/files"
	fileStore "github.com/LouisHatton/menu-link-up/internal/files/store"

	"go.uber.org/zap"
)

var _ fileStore.Reader = (*Reader)(nil)

type Reader struct {
	l          *zap.Logger
	collection string
	db         *firestore.CollectionRef
}

func New(logger *zap.Logger, collection string, client *firestore.Client) (*Reader, error) {
	r := Reader{
		l:          logger,
		collection: collection,
	}
	r.db = client.Collection(collection)
	return &r, nil
}

func (r *Reader) Get(id string) (*files.File, error) {
	logger := r.l.With(zap.String("fileId", id))

	logger.Debug("getting file doc")

	doc, err := r.db.Doc(id).Get(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error getting file: %w", err)
	}
	logger.Debug("fetched file doc")

	file := files.Empty()
	err = doc.DataTo(&file)
	if err != nil {
		return nil, fmt.Errorf("error converting response to file struct: %w", err)
	}

	return &file, nil
}

func (r *Reader) GetByUrl(id string) (*files.File, error) {
	docs, err := r.db.Where("urlId", "==", id).Documents(context.TODO()).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error fetching all documents: %w", err)
	}

	if len(docs) > 0 {
		file := files.Empty()
		doc := docs[0]
		err := doc.DataTo(&file)
		if err != nil {
			return nil, fmt.Errorf("error converting response to file: %w", err)
		} else {
			return &file, nil
		}
	} else {
		return nil, fmt.Errorf("no document found with url id")
	}
}

func (r *Reader) Many(opts query.Options, wheres ...query.Where) (*[]files.File, error) {

	q := dbFirestore.GenerateQuery(r.db.Query, opts, wheres...)

	itr := q.Documents(context.TODO())
	snapshots, err := itr.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error fetching all documents: %w", err)
	}
	docs := []files.File{}
	for i, snap := range snapshots {
		docs = append(docs, files.Empty())
		err = snap.DataTo(&docs[i])
		if err != nil {
			return nil, fmt.Errorf("error converting response to file struct: %w", err)
		}
	}

	return &docs, nil
}
