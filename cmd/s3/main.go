package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
	if err != nil {
		panic(err)
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("menulink-dev-pdf-bucket"),
		Key:    aws.String("testing.pdf"),
	})
	req.HTTPRequest.Header.Set("Content-Type", "application/pdf")
	str, err := req.Presign(15 * time.Minute)

	fmt.Println("The URL is:", str, " err:", err)

}
