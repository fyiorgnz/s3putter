package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	bucketName := os.Getenv("S3_BUCKET")

	var buf []byte
	buf, err := ioutil.ReadAll(os.Stdin)
	checkErr(err)

	auth := session.Must(session.NewSession())

	// Open Bucket
	s := s3.New(auth)
	_, err = s.HeadBucket(&s3.HeadBucketInput{Bucket: aws.String(bucketName)})
	checkErr(err)

	t := time.Now()
	fileName := formatTime(t)

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:	aws.String(fileName),
		Body:	bytes.NewReader(buf),
		ContentType:	aws.String("text/plain"),
		ACL:	aws.String("bucket-owner-full-control"),
	}
	result, err := s.PutObject(params)
	checkErr(err)
	fmt.Println(result)
	fmt.Printf("Successfully received email and saved in S3 as %s\n", fileName)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		// exit with code 75 for Postfix to know to bounce mail
		os.Exit(75)
	}
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d--%02d-%02d-%02d-%03d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond()/100000)
}
