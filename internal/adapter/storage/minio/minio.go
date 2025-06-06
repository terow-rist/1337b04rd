package minio

import (
	"1337bo4rd/internal/adapter/config"
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client     *minio.Client
	postBucket string
}

func NewMinioClient(mio *config.Minio) (*MinioClient, error) {
	// initialization of minio client
	client, err := minio.New(mio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(mio.AccessKey, mio.SecretKey, ""),
		Secure: mio.SSL,
	})
	if err != nil {
		return nil, err
	}

	// creation of posts bucket
	for attempts := 0; attempts < 3; attempts++ {
		exists, err := client.BucketExists(context.Background(), "posts")
		if err == nil && exists {
			break
		}
		if err == nil {
			err := client.MakeBucket(context.Background(), "posts", minio.MakeBucketOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to create bucket %s: %w", "posts", err)
			}
			break
		}
		time.Sleep(3 * time.Second)
	}

	return &MinioClient{
		client:     client,
		postBucket: "posts",
	}, nil
}

func (m *MinioClient) UploadImage(ctx context.Context, file multipart.File, filename, contentType string) (string, error) {
	// generate a unique object name using timestamp and filename
	objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filename)

	// upload the image to the 'posts' bucket
	_, err := m.client.PutObject(ctx, m.postBucket, objectName, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	// construct and return the URL of the uploaded image
	imageURL := fmt.Sprintf("/images/posts/%s", objectName)
	return imageURL, nil
}

func ServePostImageHandler(storage *MinioClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")
		data, contentType, err := storage.GetImage(r.Context(), "posts", filename)
		if err != nil {
			return // dumb
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func (m *MinioClient) GetImage(ctx context.Context, bucket, objectName string) ([]byte, string, error) {
	obj, err := m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object: %w", err)
	}
	defer obj.Close()

	stat, err := obj.Stat()
	if err != nil {
		return nil, "", fmt.Errorf("failed to stat object: %w", err)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(obj); err != nil {
		return nil, "", fmt.Errorf("failed to read object data: %w", err)
	}

	return buf.Bytes(), stat.ContentType, nil
}
