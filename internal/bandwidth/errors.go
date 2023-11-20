package bandwidth

import "errors"

var (
	ErrRecordNotFound               = errors.New("the bandwidth record was not found")
	ErrBytesTransferredLimitReached = errors.New("the transfer limit for this month has been reached")
	ErrUploadLimitReached           = errors.New("the upload limit for this month has been reached")
)
