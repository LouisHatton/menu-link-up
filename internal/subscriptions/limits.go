package subscriptions

import (
	"fmt"
	"strconv"
	"time"

	"github.com/LouisHatton/menu-link-up/internal/users"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stripe/stripe-go/v76"
)

func ExtractLimitsFromProduct(product *stripe.Product) (*users.BandwidthLimits, error) {
	var limits users.BandwidthLimits

	bytesTransferredLimit, err := strconv.ParseInt(product.Metadata["bytesTransferredLimit"], 10, 64)
	if err != nil {
		msg := "unable to parse bytesTransferredLimit from subscription metadata"
		return nil, fmt.Errorf(msg+": %w", err)
	}
	limits.BytesTransferredLimit = bytesTransferredLimit

	bytesUploadedLimit, err := strconv.ParseInt(product.Metadata["bytesUploadedLimit"], 10, 64)
	if err != nil {
		msg := "unable to parse bytesUploadedLimit from subscription metadata"
		return nil, fmt.Errorf(msg+": %w", err)
	}
	limits.BytesUploadedLimit = bytesUploadedLimit

	fileSizeLimit, err := strconv.ParseInt(product.Metadata["fileSizeLimit"], 10, 64)
	if err != nil {
		msg := "unable to parse fileSizeLimit from subscription metadata"
		return nil, fmt.Errorf(msg+": %w", err)
	}
	limits.FileSizeLimit = fileSizeLimit

	fileUploadLimit, err := strconv.Atoi(product.Metadata["fileUploadLimit"])
	if err != nil {
		msg := "unable to parse fileUploadLimit from subscription metadata"
		return nil, fmt.Errorf(msg+": %w", err)
	}
	limits.FileUploadLimit = fileUploadLimit

	return &limits, nil

}

func StripeTime(in int64) *time.Time {
	return aws.Time(time.Unix(in, 0))
}
