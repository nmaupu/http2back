package provider

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

var (
	_ Provider = AwsS3{}
)

type AwsS3 struct {
	Bucket, DestDir string
	// region can also be specified using AWS_REGION env var
	Region string
	// AWS creds, token is needed only for temporary creds
	// if all those vars are empty / not initialized
	// .aws/ or env vars will be used
	// aka. AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
	AwsAccessKeyId, AwsSecretAccessKey, Token string
}

// TODO : implement DestDir support
func (s AwsS3) Copy(in io.Reader, name string) string {
	// Credentials ans region handling
	var awsCreds *credentials.Credentials = nil
	if s.AwsAccessKeyId != "" && s.AwsSecretAccessKey != "" {
		awsCreds = credentials.NewStaticCredentials(s.AwsAccessKeyId, s.AwsSecretAccessKey, s.Token)
		_, err := awsCreds.Get()
		if err != nil {
			panic(fmt.Sprintf("Bad credentials - %s", err))
		}
	}

	awsConfig := aws.NewConfig()
	if s.Region != "" {
		awsConfig = awsConfig.WithRegion(s.Region)
	}
	if awsCreds != nil {
		awsConfig = awsConfig.WithCredentials(awsCreds)
	}

	// Session creation
	sess := session.Must(session.NewSession(awsConfig))
	uploader := s3manager.NewUploader(sess)

	// Upload
	destFilename := fmt.Sprintf("%s/%s", s.DestDir, name)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &s.Bucket,
		Key:    &destFilename,
		Body:   in,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to upload file to S3 - %v", err))
	}

	return fmt.Sprintf("%s", result.Location)
}

func (s AwsS3) String() string {
	return fmt.Sprintf("AwsS3 (Bucket: %s, Region: %s, Destination: %s)", s.Bucket, s.Region, s.DestDir)
}
