package storage

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
)

// GitObject represents a Git object (blob, tree, commit, or tag)
type GitObject interface {
	Type() string
	Serialize() []byte
}

// calculateHash calculates the SHA-1 hash of a Git object
func calculateHash(objType string, data []byte) string {
	h := sha1.New()
	content := fmt.Sprintf("%s %d\x00", objType, len(data))
	h.Write([]byte(content))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// ObjectReader provides an io.Reader interface for Git objects
type ObjectReader struct {
	*bytes.Reader
	hash string
	size int64
}

// NewObjectReader creates a new ObjectReader
func NewObjectReader(objType string, data []byte) *ObjectReader {
	hash := calculateHash(objType, data)
	size := int64(len(data))
	
	header := fmt.Sprintf("%s %d\x00", objType, size)
	content := append([]byte(header), data...)
	
	return &ObjectReader{
		Reader: bytes.NewReader(content),
		hash:   hash,
		size:   size,
	}
}

// Hash returns the object's hash
func (r *ObjectReader) Hash() string {
	return r.hash
}

// Size returns the object's size
func (r *ObjectReader) Size() int64 {
	return r.size
}

// ObjectWriter provides an io.Writer interface for Git objects
type ObjectWriter struct {
	buf    *bytes.Buffer
	objType string
	hash   string
}

// NewObjectWriter creates a new ObjectWriter
func NewObjectWriter(objType string) *ObjectWriter {
	return &ObjectWriter{
		buf:     bytes.NewBuffer(nil),
		objType: objType,
	}
}

func (w *ObjectWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

// Close finalizes the object and calculates its hash
func (w *ObjectWriter) Close() error {
	data := w.buf.Bytes()
	w.hash = calculateHash(w.objType, data)
	return nil
}

// Hash returns the object's hash
func (w *ObjectWriter) Hash() string {
	return w.hash
}

// Bytes returns the object's content
func (w *ObjectWriter) Bytes() []byte {
	return w.buf.Bytes()
}

// Common Git object types
const (
	BlobObject   = "blob"
	TreeObject   = "tree"
	CommitObject = "commit"
	TagObject    = "tag"
)
