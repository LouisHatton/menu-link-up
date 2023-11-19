package objectstore

import (
	"context"
	"time"

	"go.uber.org/zap/zapcore"
)

type Service interface {
	GenerateFileLocation(ctx context.Context, prefix string, suffix *string) (FileLocation, error)
	PresignedPut(ctx context.Context, location FileLocation, fileSize int, expires time.Duration) (string, error)
	PresignedGet(ctx context.Context, location FileLocation, expires time.Duration) (string, error)
	DeleteFile(ctx context.Context, location FileLocation) error
}

type FileLocation struct {
	Region string
	Bucket string
	Key    string
}

func (f FileLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("region", f.Region)
	enc.AddString("key", f.Key)
	enc.AddString("bucket", f.Bucket)
	return nil
}
