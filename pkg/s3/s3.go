package s3

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jaevor/go-nanoid"
)

type S3 struct {
	client    *awss3.Client
	presigner *awss3.PresignClient
}

func New() (*S3, error) {
	cfg, err := awscfg.LoadDefaultConfig(context.Background(),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.Get().S3.AccessKey,
				config.Get().S3.SecretKey,
				"",
			),
		),
		awscfg.WithRegion(config.Get().S3.Region),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := awss3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(client)
	return &S3{
		client:    client,
		presigner: presigner,
	}, nil
}

func (s *S3) GetClient() *awss3.Client {
	return s.client
}

func (s *S3) Upload(ctx context.Context, bucketName string, folder string, fileHeader *multipart.FileHeader) (res string, err error) {
	ctx, parentSpan := tracer.Start(ctx, "pkg.s3.Upload")
	defer func() {
		if err != nil {
			parentSpan.RecordError(err)
		}
		parentSpan.End()
	}()

	// remove whitespace
	filename := strings.ReplaceAll(fileHeader.Filename, " ", "_")

	randId, err := nanoid.Standard(15)
	if err != nil {
		return "", err
	}

	filename = randId() + "_" + filename
	filename = folder + "/" + filename
	if config.Get().Service.Name != "" {
		filename = config.Get().Service.Name + "/" + filename
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = s.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Body:          file,
		Key:           aws.String(filename),
		ContentLength: aws.Int64(fileHeader.Size),
		ContentType:   aws.String(fileHeader.Header.Get("Content-Type")),
		ACL:           types.ObjectCannedACLPublicRead,
	})
	return filename, err
}

func (s *S3) GetURL(bucketName, objectName string) string {
	baseURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucketName, config.Get().S3.Region)

	return baseURL + "/" + objectName
}

func (s *S3) GetObjectPresignedURL(ctx context.Context, bucketName, objectName string, lifetime time.Duration) (string, error) {
	req, err := s.presigner.PresignGetObject(ctx, &awss3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}, func(options *s3.PresignOptions) {
		options.Expires = lifetime
	})

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
