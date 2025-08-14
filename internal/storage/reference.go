package storage

import (
	"context"
	"fmt"
	"path"
	"strings"

	gsdk "github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
)

// ReferenceStorage implements storage for Git references in BNB Greenfield
type ReferenceStorage struct {
	client     *gsdk.GreenfieldClient
	bucketName string
	prefix     string
}

// NewReferenceStorage creates a new reference storage instance
func NewReferenceStorage(client *gsdk.GreenfieldClient, bucketName, prefix string) *ReferenceStorage {
	return &ReferenceStorage{
		client:     client,
		bucketName: bucketName,
		prefix:     prefix,
	}
}

// SetReference stores a Git reference in BNB Greenfield
func (s *ReferenceStorage) SetReference(ctx context.Context, refName, hash string) error {
	refPath := path.Join(s.prefix, "refs", refName)
	
	createObjectTx, err := s.client.CreateObject(
		ctx,
		s.bucketName,
		refPath,
		types.CreateObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to create reference: %w", err)
	}

	// Upload the reference data
	if err := s.client.UploadObject(
		ctx,
		createObjectTx,
		[]byte(hash),
		types.UploadObjectOptions{},
	); err != nil {
		return fmt.Errorf("failed to upload reference: %w", err)
	}

	return nil
}

// GetReference retrieves a Git reference from BNB Greenfield
func (s *ReferenceStorage) GetReference(ctx context.Context, refName string) (string, error) {
	refPath := path.Join(s.prefix, "refs", refName)
	
	// Download the reference
	data, err := s.client.GetObject(
		ctx,
		s.bucketName,
		refPath,
		types.GetObjectOptions{},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get reference: %w", err)
	}

	return string(data), nil
}

// DeleteReference removes a Git reference from BNB Greenfield
func (s *ReferenceStorage) DeleteReference(ctx context.Context, refName string) error {
	refPath := path.Join(s.prefix, "refs", refName)
	
	if err := s.client.DeleteObject(
		ctx,
		s.bucketName,
		refPath,
		types.DeleteObjectOptions{},
	); err != nil {
		return fmt.Errorf("failed to delete reference: %w", err)
	}

	return nil
}

// ListReferences lists all Git references in BNB Greenfield
func (s *ReferenceStorage) ListReferences(ctx context.Context) (map[string]string, error) {
	prefix := path.Join(s.prefix, "refs")
	
	objects, err := s.client.ListObjects(
		ctx,
		s.bucketName,
		types.ListObjectsOptions{
			Prefix: prefix,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list references: %w", err)
	}

	refs := make(map[string]string)
	for _, obj := range objects.Objects {
		// Extract reference name from path
		refName := strings.TrimPrefix(obj.ObjectInfo.ObjectName, prefix+"/")
		
		// Get reference value
		hash, err := s.GetReference(ctx, refName)
		if err != nil {
			return nil, fmt.Errorf("failed to get reference %s: %w", refName, err)
		}
		
		refs[refName] = hash
	}

	return refs, nil
}
