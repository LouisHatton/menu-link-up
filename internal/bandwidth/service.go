package bandwidth

import "context"

type Service interface {
	// Takes a userId and file size and attempts to log the bytes transferred
	//
	// If the user has already exceeded the monthly limit the following error is returned
	// and operation is not logged.
	//	err := bandwidth.ErrBytesTransferredLimitReached
	RecordDocumentView(ctx context.Context, userId string, fileSize int) error

	// Takes a userId and file size and attempts to log the bytes transferred
	//
	// If the user has already exceeded the monthly limit the following error is returned
	// and operation is not logged.
	//	err := bandwidth.ErrUploadLimitReached
	RecordDocumentUpload(ctx context.Context, userId string, fileSize int) error
}
