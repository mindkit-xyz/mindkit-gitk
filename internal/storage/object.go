package storage

import (
	"context"
	"fmt"
	"io"
	"path"

	gsdk "github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
)

// ObjectStorage implements storage for Git objects in BNB Greenfield
type ObjectStorage struct {
	client     *gsdk.GreenfieldClient
	bucketName string
	prefix     string
}

// NewObjectStorage creates a new object storage instance
func NewObjectStorage(client *gsdk.GreenfieldClient, bucketName, prefix string) *ObjectStorage {
	return &ObjectStorage{
		client:     client,
		bucketName: bucketName,
		prefix:     prefix,
	}
}

// Store stores a Git object in BNB Greenfield
func (s *ObjectStorage) Store(ctx context.Context, hash string, data []byte) error {
	objectPath := path.Join(s.prefix, "objects", hash[:2], hash[2:])
	
	createObjectTx, err := s.client.CreateObject(
		ctx,
		s.bucketName,
		objectPath,
		types.CreateObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to create object: %w", err)
	}

	// Upload the object data
	if err := s.client.UploadObject(
		ctx,
		createObjectTx,
		data,
		types.UploadObjectOptions{},
	); err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

// Get retrieves a Git object from BNB Greenfield
func (s *ObjectStorage) Get(ctx context.Context, hash string) ([]byte, error) {
	objectPath := path.Join(s.prefix, "objects", hash[:2], hash[2:])
	
	// Get object info
	_, err := s.client.HeadObject(
		ctx,
		s.bucketName,
		objectPath,
		types.HeadObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to head object: %w", err)
	}

	// Download the object
	data, err := s.client.GetObject(
		ctx,
		s.bucketName,
		objectPath,
		types.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return data, nil
}

// Delete removes a Git object from BNB Greenfield
func (s *ObjectStorage) Delete(ctx context.Context, hash string) error {
	objectPath := path.Join(s.prefix, "objects", hash[:2], hash[2:])
	
	if err := s.client.DeleteObject(
		ctx,
		s.bucketName,
		objectPath,
		types.DeleteObjectOptions{},
	); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// List lists all Git objects in BNB Greenfield
func (s *ObjectStorage) List(ctx context.Context) ([]string, error) {
	prefix := path.Join(s.prefix, "objects")
	
	objects, err := s.client.ListObjects(
		ctx,
		s.bucketName,
		types.ListObjectsOptions{
			Prefix: prefix,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var hashes []string
	for _, obj := range objects.Objects {
		// Extract hash from object path
		hash := path.Base(obj.ObjectInfo.ObjectName)
		hashes = append(hashes, hash)
	}

	return hashes, nil
}
