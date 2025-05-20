package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/khaledhikmat/tr-extractor/service/config"
	"github.com/khaledhikmat/tr-extractor/service/lgr"
)

type s3Service struct {
	ConfigSvc config.IService
	Ctx       context.Context
	Client    *s3.Client
}

func NewS3(ctx context.Context, cfgsvc config.IService) IService {
	s := &s3Service{
		ConfigSvc: cfgsvc,
		Ctx:       ctx,
	}
	err := s.makeS3Client(ctx)
	if err != nil {
		panic(err)
	}

	return s
}

func (svc *s3Service) Upload(filePath, folder, identifier string) (string, error) {
	bucketName := svc.ConfigSvc.GetStorageBucket()
	keyName := fmt.Sprintf("%s/%s", folder, identifier)
	lgr.Logger.Info("S3.Upload",
		slog.String("filePath", filePath),
		slog.String("folder", folder),
		slog.String("identifier", identifier),
		slog.String("bucket", bucketName),
		slog.String("key", keyName),
	)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	defer func() {
		// Delete the local file
		_ = os.Remove(filePath)
	}()

	// WARNING: if the file already exists in S3, it will be overwritten
	_, err = svc.Client.PutObject(svc.Ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	// URL must be generated from the bucket name and the folder
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, svc.ConfigSvc.GetStorageRegion(), keyName), nil
}

func (svc *s3Service) makeS3Client(ctx context.Context) error {
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(svc.ConfigSvc.GetStorageRegion()),
	)

	if err != nil {
		return err
	}

	svc.Client = s3.NewFromConfig(cfg)
	return nil
}
