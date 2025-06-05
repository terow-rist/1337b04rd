package minio

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"
	"net/http"
	"bytes"
	
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client        *minio.Client
	avatarBucket  string
	postBucket    string
	commentBucket string
	sessionBucket string 
}

func NewMinioClient(endpoint string) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Minio client: %w", err)
	}

	// Buckets to create
	buckets := []string{"avatars", "posts", "comments", "sessions"} 

	for _, bucket := range buckets {
		for attempts := 0; attempts < 3; attempts++ {
			exists, err := client.BucketExists(context.Background(), bucket)
			if err == nil && exists {
				break
			}
			if err == nil {
				err := client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
				if err != nil {
					return nil, fmt.Errorf("failed to create bucket %s: %w", bucket, err)
				}
				break
			}
			time.Sleep(3 * time.Second)
		}
	}

	return &MinioClient{
		client:        client,
		avatarBucket:  "avatars",
		postBucket:    "posts",
		commentBucket: "comments",
		sessionBucket: "sessions", 
	}, nil
}

func (m *MinioClient) UploadImage(ctx context.Context, file multipart.File, filename, contentType string) (string, error) {
	// Generate a unique object name using timestamp and filename
	objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filename)

	// Upload the image to the 'posts' bucket
	_, err := m.client.PutObject(ctx, m.postBucket, objectName, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	// Construct and return the URL of the uploaded image
	imageURL := fmt.Sprintf("/images/posts/%s", objectName)
	return imageURL, nil
}

func ServePostImageHandler(storage *MinioClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename") // Go 1.21+
		data, contentType, err := storage.GetImage(r.Context(), "posts", filename)
		if err != nil {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
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